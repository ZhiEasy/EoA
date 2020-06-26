# deploy_svr.sh
ps aux | grep svr_linux_amd64 | awk '{print $2}' | xargs kill -9
cd ~/.eoa/conf/
# wget http://jdcloud.ningdali.com/eoa/conf/app.conf
wget http://39.101.176.8/downloads/app.conf
cd ~/.eoa/
# wget http://jdcloud.ningdali.com/eoa/svr_linux_amd64
wget http://39.101.176.8/downloads/svr_linux_amd64
chmod +x svr_linux_amd64
nohup ./svr_linux_amd64 > svr.log 2>&1 &
