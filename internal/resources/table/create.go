package table

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tbl "github.com/ydb-platform/terraform-provider-ydb/internal/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

func (h *handler) Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tableResource, err := tableResourceSchemaToTableResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if tableResource == nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "got nil resource, unreachable code",
			},
		}
	}
	tableResource.Token = h.token // TODO(shmel1k@):
	db, err := tbl.CreateDBConnection(ctx, tbl.ClientParams{
		DatabaseEndpoint: tableResource.DatabaseEndpoint,
		Token:            tableResource.Token,
	})
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "failed to initialize table client",
				Detail:   err.Error(),
			},
		}
	}
	defer func() {
		_ = db.Close(ctx)
	}()

	q := PrepareCreateRequest(tableResource)
	err = db.Table().Do(ctx, func(ctx context.Context, s table.Session) (err error) {
		return s.ExecuteSchemeQuery(ctx, q)
	})
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "failed to create table",
				Detail:   err.Error(),
			},
		}
	}

	d.SetId(tableResource.Path)

	return h.Read(ctx, d, meta)
}
