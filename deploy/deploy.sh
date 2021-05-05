# 请修改wfRadius.service中的用户名及运行路径
sudo cp wfRadius.service /usr/lib/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable wfRadius.service
sudo systemctl start wfRadius.service