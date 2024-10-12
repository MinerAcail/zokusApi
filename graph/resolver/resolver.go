package graph

import (
	"github.com/mineracail/zokus/payment"
	"github.com/minio/minio-go/v7"
	"github.com/rpip/paystack-go"
	"gorm.io/gorm"
)




type Resolver struct {
	DB          *gorm.DB
	MinIOClient *minio.Client
	MoMoAPI *payment.MoMoAPI 
	
	PaystackAPI   *paystack.Client
}
