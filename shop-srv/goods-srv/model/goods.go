package model

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"shop-srv/goods-srv/global"
)

// Category 类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为 not null 或者 添加默认值
//https://zhuanlan.zhihu.com/p/73997266
//这些类型使用int32：与grpc字段类型统一
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);default:'';not null;comment:商品名称" json:"name"`
	ParentCategoryID int32       `gorm:"default:null;comment:父类目ID" json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;default:1;not null;comment:类目级别" json:"level"`
	IsTab            bool        `gorm:"default:false;not null;comment:是否在Tab栏展示" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);default:'';not null;comment:品牌名称"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:品牌LogoURL"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);default:'';not null;comment:轮播图片URL"`
	Url   string `gorm:"type:varchar(200);default:'';not null;comment:轮播图片详情URL"`
	Index int32  `gorm:"type:int;default:1;not null;comment:轮播图索引"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null;comment:是否上架"`
	ShipFree bool `gorm:"default:false;not null;comment:是否包邮"`
	IsNew    bool `gorm:"default:false;not null;comment:是否新品"`
	IsHot    bool `gorm:"default:false;not null;comment:是否热卖"`

	Name            string   `gorm:"type:varchar(50);not null;comment:商品名称"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:商品编号"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:商品点击数"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:商品销量"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:收藏数量"`
	MarketPrice     float32  `gorm:"type:decimal(10,2);not null;comment:市场价格"`
	ShopPrice       float32  `gorm:"type:decimal(10,2);not null;comment:实际价格"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:商品简介"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:商品缩略图"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:商品详情图"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:商品封面图"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
