#REST API

##REST endpoint
The Ninja Presets REST API is available on port 8101 of the master Ninja Sphere.

	http://{master-ninja-sphere}:8101

##JSON Model

###Scene

		{
		  "id" : "70a642d6-b1c1-11e4-b359-7c669d02a706",
		  "scope" : "site:a458dfe3-3a81-43cc-a118-6c42c814f4b3",
		  "slot" : 1,
		  "label" : "Preset 1",
		  "things" : [
		     {
		        "id" : "e859969e-b056-11e4-ae28-7c669d02a706",
		        "channels" : [
		           {
		              "id" : "1-6-in",
		              "state" : true
		           }
		        ]
		     }
		  ]
		}


##Methods

###GET /rest/v1/presets?scope={scope-id}
Answers a JSON array containing all the scenes in the specified scope. Scope-id is one-of 'site', or 'room:{room-id}' where room-id is the identifier of a room.

###POST /rest/v1/presets
Create a new scene using the JSON object provided in the body of the POST request. Answers the created object in the response.

####GET /rest/v1/presets/{scene-id}
Answers a JSON object containing the channel states for each thing in the scene.

####PUT /rest/v1/presets/{scene-id}
Replace the specified scene with the JSON object provided in the body of the PUT request. Answers the updated object in the response.

####DELETE /rest/v1/presets/{scene-id}
Delete the specified scene. Answers the deleted object in the response.

####POST /rest/v1/presets/{scene-id}/apply
Apply the specified scene to the scene's things.

####GET /rest/v1/presets/prototype/site
Answers a JSON object which contains a prototype scene containing the current states of each presetable thing in the site.

####GET /rest/v1/presets/prototype/room/{room-id}
Answers a JSON object which contains a prototype scene containing the current states of each presetable thing in the specified room.

##Examples

The following examples show how to use the API with 'curl' and 'jq' to achieve various tasks relating to setting and getting presets.

### Store the current state in a site-scoped preset 1 with label "from-curl"

	curl -s http://${SPHERE:-ninjasphere}:8101/rest/v1/presets/prototype/site | curl -d @- "http://${SPHERE:-ninjasphere}:8101/rest/v1/presets?slot=1&label=from-curl" | jq .

### List all existing presets

	curl -s http://${SPHERE:-ninjasphere}:8101/rest/v1/presets | jq .

### List all existing site-scoped presets

	curl -s http://${SPHERE:-ninjasphere}:8101/rest/v1/presets?scope=site | jq .

### List all existing room-scoped presets

	curl -s http://${SPHERE:-ninjasphere}:8101/rest/v1/presets?scope=room:{room-id} | jq .

### Delete all presets

	curl -s -X DELETE http://${SPHERE:-ninjasphere}:8101/rest/v1/presets | jq .
