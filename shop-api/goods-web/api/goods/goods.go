package goods

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/goods-web/forms"
	"shop-api/goods-web/proto"
	"strconv"
	"strings"

	"shop-api/goods-web/global"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func IndexGoods(ctx *gin.Context) {
	fmt.Println("主页分类商品列表")

	// 使用sentinel对这段代码进行熔断限流
	e, b := sentinel.Entry("index-goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}

	// 查询一级商品分类
	topLevelCategories := []int32{130358, 130361}
	indexGoodsList := make([]interface{}, 0)

	for _, categoryId := range topLevelCategories {
		r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
			Id: categoryId,
		})
		if err != nil {
			zap.S().Errorw("[IndexGoods] 查询 【一级商品分类】失败")
			HandleGrpcErrorToHttp(err, ctx)
			return
		}
		subCategories := make([]interface{}, 0)
		adGoodMap := make(map[string]interface{}) // 初始化广告商品
		brands := make([]interface{}, 0)          // 初始化品牌
		goods := make([]interface{}, 0)
		// 遍历二级分类
		for secondCategoryIndex, value := range r.SubCategorys {
			goodsCategoryRes, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
				Id: value.Id,
			})
			if err != nil {
				zap.S().Errorw("[IndexGoods] 查询 【三级分类ID】失败")
				HandleGrpcErrorToHttp(err, ctx)
				return
			}
			var goodsCategoriesIds []int32
			for thirdCategoryIndex, goodsCategories := range goodsCategoryRes.SubCategorys {
				if thirdCategoryIndex == 0 { // 获取第一个商品（作为广告商品）的相关品牌信息
					// 获取一级分类所包含的品牌
					categoryBrandRsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
						Id: goodsCategories.Id,
					})
					if err != nil {
						zap.S().Errorw("[IndexGoods] 查询 【一级商品分类所含品牌】失败")
						HandleGrpcErrorToHttp(err, ctx)
						return
					}
					for n, v := range categoryBrandRsp.Data {
						if secondCategoryIndex == 0 && n < 3 { // 只展示三个品牌信息
							brands = append(brands, map[string]interface{}{
								"id":    v.Id,
								"name":  v.Name,
								"image": v.Logo,
							})
						}
					}
				}
				goodsCategoriesIds = append(goodsCategoriesIds, goodsCategories.Id)
			}
			goodsRes, err := global.GoodsSrvClient.BatchGetGoodsByCates(context.Background(), &proto.BatchCategoriesIdInfo{
				CategoryId: goodsCategoriesIds,
			})
			if err != nil {
				zap.S().Errorw("[IndexGoods] 查询 【三级商品】失败")
				HandleGrpcErrorToHttp(err, ctx)
				return
			}
			for index, goodDetail := range goodsRes.Data {
				if index == 0 {
					adGoodMap["id"] = goodDetail.Id
					adGoodMap["goods_front_image"] = goodDetail.GoodsFrontImage
				}
				goods = append(goods, map[string]interface{}{
					"id":                goodDetail.Id,
					"name":              goodDetail.Name,
					"goods_front_image": goodDetail.GoodsFrontImage,
					"shop_price":        goodDetail.ShopPrice,
				})
			}

			subCategories = append(subCategories, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"category_type":   value.Level,
				"parent_category": value.ParentCategory,
				"is_tab":          value.IsTab,
			})
		}

		indexGoodsList = append(indexGoodsList, map[string]interface{}{
			"id":              categoryId,
			"name":            r.Info.Name,
			"is_tab":          r.Info.IsTab,
			"parent_category": r.Info.ParentCategory,
			"category_type":   r.Info.Level,
			"sub_cat":         subCategories,
			"brands":          brands,
			"goods":           goods,
			"ad_goods":        adGoodMap,
		})
	}
	e.Exit()

	ctx.JSON(http.StatusOK, gin.H{
		"data": indexGoodsList,
	})
}

func List(ctx *gin.Context) {
	fmt.Println("商品列表")
	//商品的列表 pmin=abc, spring cloud, go-micro
	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("pn", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//请求商品的service服务、负载均衡
	//parent, _ := ctx.Get("parentSpan")
	//opentracing.ContextWithSpan(context.Background(), parent.(opentracing.Span))

	// 使用sentinel对这段代码进行熔断限流
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}
	r, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()

	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//如何设置库存
	//TODO 商品的库存 - 分布式事务
	ctx.JSON(http.StatusOK, rsp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//查询商品库存
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("[goods.Detail] 查询【库存信息】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	goodsStocks := invRsp.Num

	rsp := map[string]interface{}{
		"id":           r.Id,
		"name":         r.Name,
		"goods_brief":  r.GoodsBrief,
		"goods_sn":     r.GoodsSn,
		"desc":         r.GoodsDesc,
		"ship_free":    r.ShipFree,
		"images":       r.Images,
		"desc_images":  r.DescImages,
		"front_image":  r.GoodsFrontImage,
		"market_price": r.MarketPrice,
		"shop_price":   r.ShopPrice,
		"stocks":       goodsStocks,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
	return
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//TODO 商品的库存
	return
}

func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
