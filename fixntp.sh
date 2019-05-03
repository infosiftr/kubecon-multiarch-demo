# Include the following in your shell session:
# ssh-agent
# ssh-add /path/to/key
for i in `seq 0 5`; do ssh root@172.23.0.10$i ntpdate -u 172.23.0.151; done

