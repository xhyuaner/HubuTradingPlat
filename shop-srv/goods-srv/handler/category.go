package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srv/goods-srv/global"
	"shop-srv/goods-srv/model"
	"shop-srv/goods-srv/proto"
)

func CategoryModel2Message(category model.Category) *proto.CategoryInfoResponse {
	return &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}
}

// GetAllCategorysList 商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*	web层封装成下面这种json数据比较麻烦，srv层直接使用gorm可以很容易封装，所以这里直接封装成JsonData返回给web层
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[
					"id":xxx,
					"name":"",
					"level":1,
					"is_tab":false,
					"sub_category":[]
				]
			},
		]
	*/
	var categorys []model.Category
	categoryListResponse := proto.CategoryListResponse{}
	result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)

	categoryListResponse.Total = int32(result.RowsAffected)
	b, _ := json.Marshal(&categorys)
	categoryListResponse.JsonData = string(b)
	// TODO:改动-1 将categoryListResponse.Data数据进行填充
	var categoryResponses []*proto.CategoryInfoResponse
	for _, category := range categorys {
		categoryResponses = append(categoryResponses, CategoryModel2Message(category))
	}
	categoryListResponse.Data = categoryResponses

	return &categoryListResponse, nil
}

// GetSubCategory 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = CategoryModel2Message(category)

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, CategoryModel2Message(subCategory))
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}
	//可能报错
	re := global.DB.Save(&category)
	fmt.Println(re)
	//cMap := map[string]interface{}{}
	//cMap["name"] = req.Name
	//cMap["level"] = req.Level
	//cMap["is_tab"] = req.IsTab
	//if req.Level != 1 {
	//	//去查询父类目是否存在
	//	cMap["parent_category_id"] = req.ParentCategory
	//}
	//tx := global.DB.Model(&model.Category{}).Create(cMap)
	//fmt.Println(tx)
	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
