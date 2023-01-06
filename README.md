# httpproxydetect
虽然这么叫，但是是一个协议高速检测工具，可以搭配zmap用,也可以搭配我的另一个项目bruteforcer利用 例如  
sudo zmap -B 50M -p 22 | ./httpproxydetect -P ssh -p 22 |./bruteforcer -P ssh -p 22 -C "whoami"> outputBrutessh.txt  
仅供合法合规的授权渗透测试，请遵守法律!
