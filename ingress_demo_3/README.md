# fixing route /backend/ correctly

```bash
$ kubectl exec -it frontend-deployment-5ccdb88c57-cdzcl -n 2048-game -- /bin/bash


# apt-get update

# apt-get install vim

# find -name index.html

# cd /usr/share/nginx/html/

# mkdir backend

# cp ./frontend/index.html ./backend/


$ kubectl apply -f ingress-demo-3.yaml
```
