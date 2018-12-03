build:
	docker build -t stephenhillier/earthworks-web . && docker push stephenhillier/earthworks-web && kubectl -n earthworks replace --force -f ../kubernetes/03-web-deploy.yaml