package index

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ydb-platform/terraform-provider-ydb/internal/helpers"
	"github.com/ydb-platform/terraform-provider-ydb/internal/resources/table/index"
	"github.com/ydb-platform/terraform-provider-ydb/sdk/terraform/auth"
)

func ResourceCreateFunc(cb auth.GetTokenCallback) helpers.TerraformCRUD {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		token, err := cb(ctx)
		if err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "failed to create token for YDB request",
					Detail:   err.Error(),
				},
			}
		}

		h := index.NewHandler(token)
		return h.Create(ctx, d, meta)
	}
}

func ResourceReadFunc(cb auth.GetTokenCallback) helpers.TerraformCRUD {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		token, err := cb(ctx)
		if err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "failed to create token for YDB request",
					Detail:   err.Error(),
				},
			}
		}

		h := index.NewHandler(token)
		return h.Read(ctx, d, meta)
	}
}

func ResourceUpdateFunc(cb auth.GetTokenCallback) helpers.TerraformCRUD {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		token, err := cb(ctx)
		if err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "failed to create token for YDB request",
					Detail:   err.Error(),
				},
			}
		}

		h := index.NewHandler(token)
		return h.Update(ctx, d, meta)
	}
}

func ResourceDeleteFunc(cb auth.GetTokenCallback) helpers.TerraformCRUD {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		token, err := cb(ctx)
		if err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "failed to create token for YDB request",
					Detail:   err.Error(),
				},
			}
		}

		h := index.NewHandler(token)
		return h.Delete(ctx, d, meta)
	}
}

func ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"table_path": {
			Type:     schema.TypeString,
			Required: true,
		},
		"connection_string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"columns": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		"cover": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
