package list

import (
	"fmt"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"k8s.io/apimachinery/pkg/util/duration"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorNamespaceList, wide bool) {
	t := dumper.Table[admin.ConnectorNamespace]{}

	t.Field("ID", func(in *admin.ConnectorNamespace) string {
		return in.Id
	})

	if wide {
		t.Field("Name", func(in *admin.ConnectorNamespace) string {
			return in.Name
		})
	}

	t.Field("ClusterID", func(in *admin.ConnectorNamespace) string {
		return in.ClusterId
	})

	t.Field("Owner", func(in *admin.ConnectorNamespace) string {
		return in.Owner
	})

	if wide {
		t.Field("TenantKind", func(in *admin.ConnectorNamespace) string {
			return string(in.Tenant.Kind)
		})
	}

	t.Field("TenantID", func(in *admin.ConnectorNamespace) string {
		return in.Tenant.Id
	})

	if wide {
		t.Field("CreatedAt", func(in *admin.ConnectorNamespace) string {
			return in.CreatedAt.Format(time.RFC3339)
		})

		t.Field("ModifiedAt", func(in *admin.ConnectorNamespace) string {
			return in.ModifiedAt.Format(time.RFC3339)
		})
	}

	if !wide {
		t.Field("Age", func(in *admin.ConnectorNamespace) string {
			age := duration.HumanDuration(time.Since(in.CreatedAt))
			if in.CreatedAt.IsZero() {
				age = ""
			}

			return age
		})
	}

	t.Rich("State", func(in *admin.ConnectorNamespace) (string, tablewriter.Colors) {
		s := string(in.Status.State)
		c := tablewriter.Colors{}

		switch s {
		case "ready":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
		case "disconnected":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
		}

		if wide && in.Expiration != "" {
			t, err := time.Parse(time.RFC3339, in.Expiration)
			if err == nil && time.Now().After(t) {
				s = s + " (*)"
			}
		}

		return s, c
	})

	t.Field("Connectors", func(in *admin.ConnectorNamespace) string {
		return fmt.Sprint(in.Status.ConnectorsDeployed)
	})

	if wide {
		t.Rich("Expiration", func(in *admin.ConnectorNamespace) (string, tablewriter.Colors) {
			s := in.Expiration
			c := tablewriter.Colors{}

			if s != "" {
				t, err := time.Parse(time.RFC3339, s)
				if err == nil && time.Now().After(t) {
					c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
				}
			}

			return s, c
		})
	}

	t.Dump(f.IOStreams.Out, items.Items)
}
