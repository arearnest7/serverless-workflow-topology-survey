## this script will help in automating the processs of creating yaml files and deploying to k8s 
echo "test identifier?"
read dir
echo "mention the name of the node this volume is being attached to"
read label 
echo "storage space ?"
read storage 
echo "path in worker node?"
read path


if [ ! -d $dir ]; then 
    mkdir $dir
fi
cd $dir

cat << EOF > initialise.yaml

apiVersion: v1
kind: Pod
metadata:
  name: $dir-initialise
spec:
  containers:
  - name: $dir-initialise
    image: busybox
    command:  ["/bin/sh", "-c", "ls && sleep 3600"]
    volumeMounts:
    - mountPath: /test
      name: $dir-local-pv
      readOnly: true
  volumes:
  - name: $dir-local-pv
    hostPath:
      path: $path
      type: Directory


EOF



echo "do you want to make any deployments? 0 for no 1 for yes"
read response 



## deployments -> ask for image, size while claiming the volume and subsequently make file for each deployment 

i=1
while [ "$response" -eq '1' ] ; 
do 
echo "image?"
read image 
echo "size in pv?"
read size_pvc
echo $i

cat << EOF > deployment-$i.yaml

apiVersion: v1
kind: Pod
metadata:
  name: $dir-deployment-$i-initialise
spec:
  containers:
  - name: $dir-initialise
    image: busybox
    command:  ["/bin/sh", "-c", "ls && sleep 3600"]
    volumeMounts:
    - mountPath: /test
      name: $dir-local-pv
      readOnly: true
  volumes:
  - name: $dir-local-pv
    hostPath:
      path:  $path/deployment$i
      type: DirectoryOrCreate 

--- 


kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: $dir-deployment-$i-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer


---

apiVersion: v1
kind: PersistentVolume
metadata:
  name: $dir-deployment-$i-pv
spec:
  capacity:
    storage: $storage
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: $dir-deployment-$i-storage
  local:
    path: $path/deployment$i
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - $label

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: $dir-deployment-$i-pvc
spec:
  accessModes:
  - ReadWriteOnce
  storageClassName: $dir-deployment-$i-storage
  resources:
    requests:
      storage: $size_pvc 

---


apiVersion: apps/v1
kind: Deployment
metadata:
  name: $dir-deployment-$i
  labels:
    app: $dir-deployment-$i
spec:
  replicas: 1
  selector:
    matchLabels:
      app: $dir-deployment-$i
  template:
    metadata:
      labels:
        app: $dir-deployment-$i
    spec:
      nodeSelector:
        kubernetes.io/hostname: $label
      volumes:
        - name: $dir-deployment-$i-local-persistent-storage
          persistentVolumeClaim:
            claimName: $dir-deployment-$i-pvc
        - name: $dir-deployment-$i-local-pv
          hostPath:
            path: $path/deployment$i
            type: DirectoryOrCreate 
      initContainers:
      - name: $dir-deployment-$i-init
        image: $image
        volumeMounts:
          - name: $dir-deployment-$i-local-persistent-storage
            mountPath: /mnt
        securityContext:
          runAsUser: 0  
        command: ["sh", "-c", "cp -r /workspace/* /mnt/"]
      containers:
      - name: $dir-deployment-$i
        image: $image
        volumeMounts:
          - name: $dir-deployment-$i-local-persistent-storage
            mountPath: /mnt
        securityContext:
          runAsUser: 0  
        ports:
        - containerPort: 8080
EOF
 ((i++))
echo "do you have more deployments ?"
read response 
done



echo "do you want to apply the files ? 0 or 1"
read apply

if [ $apply -eq 1 ]; then
for file in ./*.yaml; do
echo "Applying $file..."
kubectl apply -f "$file"
done
fi