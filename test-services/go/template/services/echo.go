package services

import (
	"the-service/spec/echo"
	"the-service/spec/empty"
	"the-service/spec/models"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EchoService struct{}

func (service *EchoService) EchoBodyString(body string) (*string, error) {
	return &body, nil
}
func (service *EchoService) EchoBodyModel(body *models.Message) (*models.Message, error) {
	return body, nil
}
func (service *EchoService) EchoBodyArray(body *[]string) (*[]string, error) {
	return body, nil
}
func (service *EchoService) EchoBodyMap(body *map[string]string) (*map[string]string, error) {
	return body, nil
}
func (service *EchoService) EchoFormData(intParam int, longParam int64, floatParam float32, doubleParam float64, decimalParam decimal.Decimal, boolParam bool, stringParam string, stringOptParam *string, stringDefaultedParam string, stringArrayParam []string, uuidParam uuid.UUID, dateParam civil.Date, dateArrayParam []civil.Date, datetimeParam civil.DateTime, enumParam models.Choice) (*models.FormParameters, error) {
	return &models.FormParameters{IntField: intParam, LongField: longParam, FloatField: floatParam, DoubleField: doubleParam, DecimalField: decimalParam, BoolField: boolParam, StringField: stringParam, StringOptField: stringOptParam, StringDefaultedField: stringDefaultedParam, StringArrayField: stringArrayParam, UuidField: uuidParam, DateField: dateParam, DateArrayField: dateArrayParam, DatetimeField: datetimeParam, EnumField: enumParam}, nil
}
func (service *EchoService) EchoFormUrlencoded(intParam int, longParam int64, floatParam float32, doubleParam float64, decimalParam decimal.Decimal, boolParam bool, stringParam string, stringOptParam *string, stringDefaultedParam string, stringArrayParam []string, uuidParam uuid.UUID, dateParam civil.Date, dateArrayParam []civil.Date, datetimeParam civil.DateTime, enumParam models.Choice) (*models.FormParameters, error) {
	return &models.FormParameters{IntField: intParam, LongField: longParam, FloatField: floatParam, DoubleField: doubleParam, DecimalField: decimalParam, BoolField: boolParam, StringField: stringParam, StringOptField: stringOptParam, StringDefaultedField: stringDefaultedParam, StringArrayField: stringArrayParam, UuidField: uuidParam, DateField: dateParam, DateArrayField: dateArrayParam, DatetimeField: datetimeParam, EnumField: enumParam}, nil
}
func (service *EchoService) EchoQuery(intQuery int, longQuery int64, floatQuery float32, doubleQuery float64, decimalQuery decimal.Decimal, boolQuery bool, stringQuery string, stringOptQuery *string, stringDefaultedQuery string, stringArrayQuery []string, uuidQuery uuid.UUID, dateQuery civil.Date, dateArrayQuery []civil.Date, datetimeQuery civil.DateTime, enumQuery models.Choice) (*models.Parameters, error) {
	return &models.Parameters{IntField: intQuery, LongField: longQuery, FloatField: floatQuery, DoubleField: doubleQuery, DecimalField: decimalQuery, BoolField: boolQuery, StringField: stringQuery, StringOptField: stringOptQuery, StringDefaultedField: stringDefaultedQuery, StringArrayField: stringArrayQuery, UuidField: uuidQuery, DateField: dateQuery, DateArrayField: dateArrayQuery, DatetimeField: datetimeQuery, EnumField: enumQuery}, nil
}
func (service *EchoService) EchoHeader(intHeader int, longHeader int64, floatHeader float32, doubleHeader float64, decimalHeader decimal.Decimal, boolHeader bool, stringHeader string, stringOptHeader *string, stringDefaultedHeader string, stringArrayHeader []string, uuidHeader uuid.UUID, dateHeader civil.Date, dateArrayHeader []civil.Date, datetimeHeader civil.DateTime, enumHeader models.Choice) (*models.Parameters, error) {
	return &models.Parameters{IntField: intHeader, LongField: longHeader, FloatField: floatHeader, DoubleField: doubleHeader, DecimalField: decimalHeader, BoolField: boolHeader, StringField: stringHeader, StringOptField: stringOptHeader, StringDefaultedField: stringDefaultedHeader, StringArrayField: stringArrayHeader, UuidField: uuidHeader, DateField: dateHeader, DateArrayField: dateArrayHeader, DatetimeField: datetimeHeader, EnumField: enumHeader}, nil
}
func (service *EchoService) EchoUrlParams(intUrl int, longUrl int64, floatUrl float32, doubleUrl float64, decimalUrl decimal.Decimal, boolUrl bool, stringUrl string, uuidUrl uuid.UUID, dateUrl civil.Date, datetimeUrl civil.DateTime, enumUrl models.Choice) (*models.UrlParameters, error) {
	return &models.UrlParameters{IntField: intUrl, LongField: longUrl, FloatField: floatUrl, DoubleField: doubleUrl, DecimalField: decimalUrl, BoolField: boolUrl, StringField: stringUrl, UuidField: uuidUrl, DateField: dateUrl, DatetimeField: datetimeUrl, EnumField: enumUrl}, nil
}
func (service *EchoService) EchoEverything(body *models.Message, floatQuery float32, boolQuery bool, uuidHeader uuid.UUID, datetimeHeader civil.DateTime, dateUrl civil.Date, decimalUrl decimal.Decimal) (*echo.EchoEverythingResponse, error) {
	return &echo.EchoEverythingResponse{Ok: &models.Everything{BodyField: *body, FloatQuery: floatQuery, BoolQuery: boolQuery, UuidHeader: uuidHeader, DatetimeHeader: datetimeHeader, DateUrl: dateUrl, DecimalUrl: decimalUrl}}, nil
}
func (service *EchoService) SameOperationName() (*echo.SameOperationNameResponse, error) {
	return &echo.SameOperationNameResponse{Ok: &empty.Value}, nil
}
