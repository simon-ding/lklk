package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"github.com/simon-ding/lklk/models"
	"io"
	"reflect"
	"strconv"
)

func Wrap(handler APIHandler) gin.HandlerFunc {

	return func(c *gin.Context) {
		v := reflect.ValueOf(handler)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		fields := deepFields(v.Type())
		for _, field := range fields {
			f := v.FieldByName(field.Name)
			if !f.CanSet() {
				continue
			}

			f = indirectValue(f)

			if _, ok := f.Interface().(gin.Context); ok {
				f.Set(reflect.ValueOf(*c))
			} else if _, ok := f.Interface().(gorm.DB); ok {
				f.Set(reflect.ValueOf(models.GetDB(c)))
			}

			ginTag := field.Tag.Get("http")
			if ginTag == "" {
				continue
			}

			switch ginTag {
			case "json":
				if c.ContentType() == "application/json" {
					if err := bindingRecursive(f, c, binding.JSON); err != nil {
						c.JSON(200, ErrorReturn(err))
						return
					}
				}
			case "query":
				if err := bindingRecursive(f, c, binding.Query); err != nil {
					c.JSON(200, ErrorReturn(err))
					return
				}

			case "form":
				if c.ContentType() == "multipart/form-data" {
					if err := bindingRecursive(f, c, binding.FormMultipart); err != nil {
						c.JSON(200, ErrorReturn(err))
						return
					}
				}
			default:
				p := c.Param(ginTag)
				if p == "" {
					continue
				}
				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					i, err := strconv.Atoi(p)
					if err != nil {
						c.JSON(200, ErrorReturn(fmt.Errorf("parameter %s not valid", p)))
						return
					}
					f.SetInt(int64(i))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					i, err := strconv.Atoi(p)
					if err != nil {
						c.JSON(200, ErrorReturn(fmt.Errorf("parameter %s not valid", p)))
						return
					}
					f.SetUint(uint64(i))
				case reflect.String:
					f.Set(reflect.ValueOf(p))
				}
			}
			//handle limit
			if q, ok := f.Interface().(Limiter); ok {
				if q.GetLimit() == 0 {
					q.SetLimit(-1)
					f.Set(reflect.ValueOf(q))
				}
			} else if q, ok := f.Addr().Interface().(Limiter); ok {
				if q.GetLimit() == 0 {
					q.SetLimit(-1)
					f.Set(reflect.ValueOf(q).Elem())
				}
			}
		}
		if validator, ok := handler.(Validator); ok {
			if err := validator.Validate(); err != nil {
				c.JSON(200, ErrorReturn(err))
				return
			}
		}

		resp, err := handler.Handle()
		if err != nil {
			c.JSON(200, ErrorReturn(err))
			return
		}
		c.JSON(200, resp)

	}
}

func bindingRecursive(v reflect.Value, c *gin.Context, binding binding.Binding) error {
	v = indirectValue(v)
	if err := c.ShouldBindWith(v.Addr().Interface(), binding); err != nil && err != io.EOF {
		return err
	}
	anonymousFields := anonymousFields(v.Type())
	for _, af := range anonymousFields {
		nestField := v.FieldByName(af.Name)
		if err := bindingRecursive(nestField, c, binding); err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirect(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func anonymousFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if reflectType = indirect(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, v)
			}
		}

	}
	return fields
}

func indirect(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		if reflectValue.Elem().Kind() != reflect.Ptr {
			reflectValue.Set(reflect.New(reflectValue.Type().Elem()))
		}
		reflectValue = reflectValue.Elem()
	}
	return reflectValue

}

type APIHandler interface {
	Handle() (*Response, error) //具体api逻辑
}

type Validator interface {
	Validate() error               //权限验证相关
}

type Limiter interface {
	SetLimit(int)
	GetLimit() int
}

