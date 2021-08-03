package tf_acceptance

pass {
    true
}

violations[result] {
	result = {
		"pass": true,
		"id": "valid",
		"name": "name",
		"description": "same as minimal.rego, but with a different description",
		"message": "message",
	}
}
