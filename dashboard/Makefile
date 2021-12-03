init_dashboard:
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.4.0/aio/deploy/recommended.yaml

k8s_token:
	@kubectl -n kubernetes-dashboard get secret \
		$$(kubectl -n kubernetes-dashboard get sa/admin-user -o jsonpath="{.secrets[0].name}") \
		-o go-template="{{.data.token | base64decode}}"
