package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotevents"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAwsIotDetectorModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsIotDetectorCreate,
		Read:   resourceAwsIotDetectorRead,
		Update: resourceAwsIotDetectorUpdate,
		Delete: resourceAwsIotDetectorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"definition": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initial_state_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 200),
						},
						"states": {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
						},
					},
				}},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateArn,
			},
		},
	}
}

func resourceAwsIotDetectorCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ioteventsconn

	detectorName := d.Get("name").(string)
	detectorDefinition := d.Get("definition").(map[string]interface{})

	// How to convert list of structures to appropriate format usign aws. package
	detectorDefinitionParams := &iotevents.DetectorModelDefinition{
		InitialStateName: aws.String(detectorDefinition["initial_state_name"].(string)),
		States:           expandStringList(detectorDefinition["states"].([]interface{})),
	}

	roleArn := d.Get("role_arn").(string)

	params := &iotevents.CreateDetectorModelInput{
		DetectorModelName:       aws.String(detectorName),
		DetectorModelDefinition: detectorDefinitionParams,
		RoleArn:                 aws.String(roleArn),
	}

	if v, ok := d.GetOk("description"); ok {
		params.DetectorModelDescription = aws.String(v.(string))
	}

	if v, ok := d.GetOk("key"); ok {
		params.Key = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating IoT Model Detector: %s", params)
	_, err := conn.CreateDetectorModel(params)

	if err != nil {
		return err
	}

	d.SetId(detectorName)

	return resourceAwsIotDetectorRead(d, meta)
}

func resourceAwsIotDetectorRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ioteventsconn

	params := &iotevents.DescribeDetectorModelInput{
		DetectorModelName: aws.String(d.Id()),
	}
	log.Printf("[DEBUG] Reading IoT Events Detector Model: %s", params)
	out, err := conn.DescribeDetectorModel(params)

	if err != nil {
		return err
	}

	return nil
}

func resourceAwsIotDetectorUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ioteventsconn

	detectorName := d.Get("name").(string)
	detectorDefinition := d.Get("definition").(map[string]interface{})

	detectorDefinitionParams := &iotevents.DetectorModelDefinition{
		InitialStateName: aws.String(detectorDefinition["initial_state_name"].(string)),
		States:           aws.String(detectorDefinition["initial_state_name"].([]string)),
	}
	roleArn := d.Get("role_arn").(string)

	params := &iotevents.UpdateDetectorModelInput{
		DetectorModelName:       aws.String(detectorName),
		DetectorModelDefinition: detectorDefinitionParams,
		RoleArn:                 aws.String(roleArn),
	}

	if v, ok := d.GetOk("description"); ok {
		params.DetectorModelDescription = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Updating IoT Events Detector Model: %s", params)
	_, err := conn.UpdateDetectorModel(params)

	if err != nil {
		return err
	}

	return resourceAwsIotDetectorRead(d, meta)
}

func resourceAwsIotDetectorDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ioteventsconn

	params := &iotevents.DeleteDetectorModelInput{
		DetectorModelName: aws.String(d.Id()),
	}
	log.Printf("[DEBUG] Deleting IoT Events Detector Model: %s", params)
	_, err := conn.DeleteDetectorModel(params)

	return err
}
