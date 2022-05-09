#!/bin/bash

#give permission for everything in the express-app directory
sudo chmod -R 777 /home/ec2-user/golang-res

#navigate into our working directory where we have all our github files
cd /home/ec2-user/golang-res

#add npm and node to path
export NVM_DIR="$HOME/.nvm"	
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # loads nvm	
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # loads nvm bash_completion (node is in path now)

#install node modules
npm install
npm install pm2
#start our node app in the background
# node app.js > app.out.log 2> app.err.log < /dev/null & 
pm2 restart app.js --name "golang-res" -f
