// Code generated by iacg; DO NOT EDIT.
package teo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionCreate,
		Read:   resourceTencentCloudTeoFunctionRead,
		Update: resourceTencentCloudTeoFunctionUpdate,
		Delete: resourceTencentCloudTeoFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},

			"function_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Function.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Function description, maximum support of 60 characters.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Function content, currently only supports JavaScript code, with a maximum size of 5MB.",
			},

			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default domain name for the function.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId     string
		functionId string
	)
	var (
		request  = teov20220901.NewCreateFunctionRequest()
		response = teov20220901.NewCreateFunctionResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo function failed, reason:%+v", logId, err)
		return err
	}

	functionId = *response.Response.FunctionId

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"false"},
		Refresh:    resourceTeoFunctionCreateStateRefreshFunc_0_0(ctx, zoneId, functionId),
		Target:     []string{"true"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}
	d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionRead(d, meta)
}

func resourceTencentCloudTeoFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoFunctionById(ctx, zoneId, functionId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_function` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.FunctionId != nil {
		_ = d.Set("function_id", respData.FunctionId)
		functionId = *respData.FunctionId
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	needChange := false
	mutableArgs := []string{"remark", "content"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyFunctionRequest()

		request.ZoneId = helper.String(zoneId)

		request.FunctionId = helper.String(functionId)

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo function failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionRead(d, meta)
}

func resourceTencentCloudTeoFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	var (
		request  = teov20220901.NewDeleteFunctionRequest()
		response = teov20220901.NewDeleteFunctionResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.FunctionId = helper.String(functionId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo function failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}

func resourceTeoFunctionCreateStateRefreshFunc_0_0(ctx context.Context, zoneId string, functionId string) resource.StateRefreshFunc {
	var req *teov20220901.DescribeFunctionsRequest
	t := template.New("gotpl")
	var tplObj *template.Template
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = teov20220901.NewDescribeFunctionsRequest()
			req.ZoneId = helper.String(zoneId)

			req.FunctionIds = []*string{helper.String(functionId)}

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionsWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		if tplObj == nil {
			tplObj, err = t.Parse("{{ if .Functions }}{{ $firstFunction := index .Functions 0 }}{{ if $firstFunction.Domain }}{{ true }}{{ else }}{{ false }}{{ end }}{{ end }}")
			if err != nil {
				return resp.Response, "", fmt.Errorf("parse state go-template error: %w", err)
			}
		}
		stream := new(bytes.Buffer)
		if err := tplObj.Execute(stream, resp.Response); err != nil {
			return resp.Response, "", err
		}
		stateBytes, err := io.ReadAll(stream)
		if err != nil {
			return resp.Response, "", err
		}
		state := string(stateBytes)
		return resp.Response, state, nil
	}
}
