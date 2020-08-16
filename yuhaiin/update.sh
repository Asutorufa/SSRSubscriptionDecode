wget https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf -O apple.china.conf
wget https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf -O google.china.conf
wget https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf -O accelerated-domains.china.conf
wget https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts -O ad_wars_hosts
wget "https://pgl.yoyo.org/adservers/serverlist.php?hostformat=adblock&showintro=0&mimetype=plaintext" -O pglyoyo.txt
wget https://logroid.github.io/adaway-hosts/hosts_no_white.txt -O ja_ad

cat accelerated-domains.china.conf | sed $'s/\r$//' | sed 's/server=\//\*\./g' | sed 's/\/114\.114\.114\.114/ DIRECT/g' > yuhaiin.conf
cat google.china.conf | sed $'s/\r$//' | sed 's/server=\//\*\./g' | sed 's/\/114\.114\.114\.114/ DIRECT/g' >> yuhaiin.conf
cat apple.china.conf | sed $'s/\r$//' | sed 's/server=\//\*\./g' | sed 's/\/114\.114\.114\.114/ DIRECT/g' >> yuhaiin.conf
cat ../cn/cn.acl | sed 's/$/ DIRECT/g' >> yuhaiin.conf
cat ../common/lan.acl | sed 's/$/ DIRECT/g' >> yuhaiin.conf
cat abroad.conf | sed 's/^/\*\./g' | sed 's/$/ PROXY/g' >> yuhaiin.conf
cat abroad_ip.conf | sed 's/$/ PROXY/g' >> yuhaiin.conf

# AD
cat yuhaiin.conf > yuhaiin_ad.conf
cat pglyoyo.txt | sed $'s/\r$//' | sed '/*/!d' |sed 's/$/ BLOCK/g' >> yuhaiin_ad.conf
cat ja_ad | sed $'s/\r$//' | sed 's/127.0.0.1 //g' | sed '/#/'d  | sed '/^\s*$/d'  |sed 's/$/ BLOCK/g' >> yuhaiin_ad.conf
cat ad_wars_hosts | sed $'s/\r$//' | sed 's/127.0.0.1 //g' | sed '/#/'d | sed '1,2d' | sed '/^ *$/d' |sed 's/$/ BLOCK/g' >> yuhaiin_ad.conf
cat cn_ad.conf  | sed $'s/\r$//' | sed '/#/'d  | sed '/^\s*$/d'  |sed 's/$/ BLOCK/g' >> yuhaiin_ad.conf

# custom
cat custom.conf | sed '/#/'d  >> yuhaiin.conf
cat custom.conf | sed '/#/'d >> yuhaiin_ad.conf

# private
cat yuhaiin_ad.conf > yuhaiin_my.conf
cat private.conf >> yuhaiin_my.conf

# DELETE
rm apple.china.conf google.china.conf accelerated-domains.china.conf ad_wars_hosts pglyoyo.txt ja_ad
