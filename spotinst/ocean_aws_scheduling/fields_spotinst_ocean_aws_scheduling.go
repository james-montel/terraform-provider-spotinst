package ocean_aws_scheduling

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[ScheduledTask] = commons.NewGenericField(
		commons.OceanAWSScheduling,
		ScheduledTask,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Tasks): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(tasksIsEnabled): {
									Type:     schema.TypeBool,
									Required: true,
								},

								string(TaskType): {
									Type:     schema.TypeString,
									Required: true,
								},

								string(CronExpression): {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					string(ShutdownHours): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(IsEnabled): {
									Type:     schema.TypeBool,
									Optional: true,
								},

								string(TimeWindows): {
									Type:     schema.TypeList,
									Required: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if scheduling, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					elastigroup.SetScheduling(scheduling)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			var scheduling *aws.Scheduling = nil
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if interfaces, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					scheduling = interfaces
				}
			}
			elastigroup.SetScheduling(scheduling)
			return nil
		},
		nil,
	)

}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//func flattenAzureGroupScheduledTasks(tasks []*aws.Tasks) []interface{} {
//	result := make([]interface{}, 0, len(tasks))
//	for _, t := range tasks {
//		m := make(map[string]interface{})
//		m[string(IsEnabled)] = spotinst.BoolValue(t.IsEnabled)
//		m[string(TaskType)] = spotinst.StringValue(t.TaskType)
//		m[string(CronExpression)] = spotinst.StringValue(t.CronExpression)
//		result = append(result, m)
//	}
//	return result
//}

func expandScheduledTasks(data interface{}) (*aws.Scheduling, error) {
	scheduling := &aws.Scheduling{}
	list := data.([]interface{})
	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(Tasks)]; ok {
			tasks, err := expandtasks(v)
			if err != nil {
				return nil, err
			}
			if tasks != nil {

				scheduling.SetTasks(tasks)
			}
		}

		if v, ok := m[string(ShutdownHours)]; ok {
			shutdownHours, err := expandShutdownHours(v)
			if err != nil {
				return nil, err
			}
			if shutdownHours != nil {
				if scheduling.ShutdownHours == nil {
					scheduling.SetShutdownHours(&aws.ShutdownHours{})
				}
				scheduling.SetShutdownHours(shutdownHours)
			}
		}
	}

	return scheduling, nil
}

func expandShutdownHours(data interface{}) (*aws.ShutdownHours, error) {
	if list := data.([]interface{}); list != nil && len(list) > 0 && list[0] != nil {
		runner := &aws.ShutdownHours{}
		m := list[0].(map[string]interface{})

		var isEnabled = spotinst.Bool(false)
		if v, ok := m[string(IsEnabled)].(bool); ok {
			isEnabled = spotinst.Bool(v)
		}
		runner.SetIsEnabled(isEnabled)

		var timeWindows []string = nil
		if v, ok := m[string(TimeWindows)].([]interface{}); ok && len(v) > 0 {
			timeWindowList := make([]string, 0, len(v))
			for _, timeWindow := range v {
				if v, ok := timeWindow.(string); ok && len(v) > 0 {
					timeWindowList = append(timeWindowList, v)
				}
			}
			timeWindows = timeWindowList
		}
		runner.SetTimeWindows(timeWindows)

		return runner, nil
	}

	return nil, nil
}

func expandtasks(data interface{}) ([]*aws.Task, error) {
	list := data.(*schema.Set).List()
	tasks := make([]*aws.Task, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		task := &aws.Task{}

		if v, ok := m[string(IsEnabled)].(bool); ok {
			task.SetIsEnabled(spotinst.Bool(v))
		}

		if v, ok := m[string(TaskType)].(string); ok && v != "" {
			task.SetTaskType(spotinst.String(v))
		}

		if v, ok := m[string(CronExpression)].(string); ok && v != "" {
			task.SetCronExpression(spotinst.String(v))
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
