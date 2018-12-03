build:
	docker build -t stephenhillier/earthworks-api . && docker push stephenhillier/earthworks-api && kubectl -n earthworks replace --force -f ../kubernetes/02-api-deploy.yaml
