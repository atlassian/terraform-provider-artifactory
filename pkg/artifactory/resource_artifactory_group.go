package artifactory

import (
	"context"
	"fmt"
	"net/http"

	"github.com/atlassian/go-artifactory/v3/artifactory"
	v1 "github.com/atlassian/go-artifactory/v3/artifactory/v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArtifactoryGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_join": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"admin_privileges": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"realm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateLowerCase,
			},
			"realm_attributes": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func unmarshalGroup(s *schema.ResourceData) (*v1.Group, error) {
	d := &ResourceData{s}

	group := new(v1.Group)

	group.Name = d.getStringRef("name", false)
	group.Description = d.getStringRef("description", false)
	group.AutoJoin = d.getBoolRef("auto_join", false)
	group.AdminPrivileges = d.getBoolRef("admin_privileges", false)
	group.Realm = d.getStringRef("realm", false)
	group.RealmAttributes = d.getStringRef("realm_attributes", false)

	// Validator
	if group.AdminPrivileges != nil && group.AutoJoin != nil && *group.AdminPrivileges && *group.AutoJoin {
		return nil, fmt.Errorf("error: auto_join cannot be true if admin_privileges is true")
	}

	return group, nil
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)

	group, err := unmarshalGroup(d)
	if err != nil {
		return err
	}

	if resp, err := c.V1.Security.CreateOrReplaceGroup(context.Background(), *group.Name, group); err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("request failed: %s %s", resp.Status(), resp)
	}

	d.SetId(*group.Name)
	if err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if ok, err := resourceGroupExists(d, m); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error reading group: %s", err))
		} else if !ok {
			return resource.RetryableError(fmt.Errorf("expected group to be created, but currently not found"))
		}
		return nil
	}); err != nil {
		return err
	}
	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	rt := m.(*artifactory.Artifactory)

	group, resp, err := rt.V1.Security.GetGroup(context.Background(), d.Id())

	if err != nil {
		return err
	} else if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("request failed: %s %s", resp.Status(), resp)
	}

	echain, errs := chainError()
	echain(d.Set("name", group.Name))
	echain(d.Set("description", group.Description))
	echain(d.Set("auto_join", group.AutoJoin))
	echain(d.Set("admin_privileges", group.AdminPrivileges))
	echain(d.Set("realm", group.Realm))
	echain(d.Set("realm_attributes", group.RealmAttributes))
	if len(errs) != 0 {
		return fmt.Errorf("failed to marshal group: %v", errs)
	}
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)
	group, err := unmarshalGroup(d)
	if err != nil {
		return err
	}

	resp, err := c.V1.Security.UpdateGroup(context.Background(), d.Id(), group)
	if err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("request failed: %s %s", resp.Status(), resp)
	}

	d.SetId(*group.Name)
	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)
	group, err := unmarshalGroup(d)
	if err != nil {
		return err
	}

	if resp, err := c.V1.Security.DeleteGroup(context.Background(), *group.Name); err != nil {
		return err
	} else if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("request failed: %s %s", resp.Status(), resp)
	}
	return nil
}

func resourceGroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	rt := m.(*artifactory.Artifactory)

	_, resp, err := rt.V1.Security.GetGroup(context.Background(), d.Id())
	if err != nil {
		return false, err
	} else if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("request failed: %s %s", resp.Status(), resp)
	}

	return true, nil
}
