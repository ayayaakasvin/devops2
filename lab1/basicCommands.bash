ssh -NT -R 2222:localhost:22 ayayasvin:10.0.2.2 # in vm
ssh -p 2222 ayayasvin@localhost # in host and make sure that ssh.service runs!

ansible-playbook -i ./training/inventory/hosts.yml ./training/site.yml --ask-pass --ask-become-pass # to run the ansible
ssh -p 2222 -L 8080:localhost:8080 ayayasvin@localhost # forward keycloak