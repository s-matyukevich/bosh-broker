#bash

cat >> ~/.bashrc <<\EOF
export GEM_HOME=~/rubygems/gems
export PATH=$PATH:$HOME/rubygems/gems/bin
EOF

source ~/.bashrc
gem install bosh_cli --no-ri --no-rdoc --no-user-install

PATH=$PATH bosh-broker
