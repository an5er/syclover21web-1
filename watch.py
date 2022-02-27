import getopt
import re
import ssl
import socket
import json
import requests
import random
import time
from urllib.parse import urlparse
import sys
from concurrent.futures import ProcessPoolExecutor


rep_head = ""
rep_data = ""
req_head = ""
baseurl = ""
txt_model = False
file_path = ""
url_list = []
my_headers = [
        "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:30.0) Gecko/20100101 Firefox/30.0",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/537.75.14",
        "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)",
        'Mozilla/5.0 (Windows; U; Windows NT 5.1; it; rv:1.8.1.11) Gecko/20071127 Firefox/2.0.0.11',
        'Opera/9.25 (Windows NT 5.1; U; en)',
        'Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)',
        'Mozilla/5.0 (compatible; Konqueror/3.5; Linux) KHTML/3.5.5 (like Gecko) (Kubuntu)',
        'Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.8.0.12) Gecko/20070731 Ubuntu/dapper-security Firefox/1.5.0.12',
        'Lynx/2.8.5rel.1 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/1.2.9',
        "Mozilla/5.0 (X11; Linux i686) AppleWebKit/535.7 (KHTML, like Gecko) Ubuntu/11.04 Chromium/16.0.912.77 Chrome/16.0.912.77 Safari/535.7",
        "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:10.0) Gecko/20100101 Firefox/10.0 "
    ]
headers = {'User-Agent': random.choice(my_headers),'Connection': 'close','Accept':'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8','Accept-Encoding':'gzip, deflate','Accept-Language':'zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2','Cookie':'PHPSESSID=dvftgg9eq768afdiuqthrukn31'}


def main():
    global baseurl
    global txt_model
    global file_path
    global url_list

    opts, args = getopt.getopt(sys.argv[1:], '-h-u:-r:' , ['help'])
    for opt_name, opt_value in opts:
        if opt_name in ('-h', '--help'):
            print("[*] Help info:")
            print("[*] python watch.py -u url -h help -r urltxt ")
            sys.exit()
        if opt_name in ('-u'):
            baseurl = opt_value
        if opt_name in ('-r'):
            txt_model = True
            file_path = opt_value
    if txt_model == False:
        run(baseurl)
    else:
        f = open(file_path,'r')
        for line in f.readlines():
            line = line.strip('\n')
            url_list.append(line)
        f.close()
        with ProcessPoolExecutor() as pool:
            results = pool.map(run, url_list)
            print("正在并行执行")
    print("done")
    sys.exit()

def run(url):
    get_request(url)
    with open('fofa.json', 'r', encoding='utf-8') as f:
        fofa_dic = json.load(f)
        for i in range(len(fofa_dic)):
            for j in range(len(fofa_dic[i]['rules'])):
                print("正在检测是否为" + fofa_dic[i]["product"])
                a = []
                if (len(fofa_dic[i]['rules'][j])>=2):
                    for q in range(len(fofa_dic[i]['rules'][j])):
                        a.append(check_match(fofa_dic[i]['rules'][j][q]["match"],fofa_dic[i]['rules'][j][q]["content"],url))
                    if False not in a:
                         if True in a:
                            resulte_write(url,fofa_dic[i]["product"])

                else:
                    if (check_match(fofa_dic[i]['rules'][j][0]["match"],fofa_dic[i]['rules'][j][0]["content"],url)):
                        resulte_write(url,fofa_dic[i]["product"])

def resulte_write(url,value):
    f = open(time.strftime("%m-%d", time.localtime())+".txt",'a')
    f.write(str(url)+" cms :"+str(value)+"\n")
    f.close()
#
#   sys.exit()
#

def check_match(match,value,url):
    if match == "body_contains":
        return get_body(value)

    elif match == "protocol_contains":
        return False
    elif match == "title_contains":
        return get_title(value)
    elif match == "banner_contains":
        return get_banner(value)
    elif match == "header_contains":
        return get_head(value)
    elif match ==  "port_contains":
        return get_port(url,value)
    elif match == "server":
        return get_head(value)
    elif match == "title":
        return get_title(value)
    elif match == "cert_contains":
        return False
#        get_cert(url,value)
    elif match == "server_contains":
        return get_server(value)
    elif match == "protocol":
        return False
    else:
        print(match+"匹配失败")

def isIP(str):
    p = re.compile('^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$')
    if p.match(str):
        return True
    else:
        return False

def get_port(host,value):
    if isIP(host):
        _url = urlparse(host)
        port = _url.port
        if str(port) == value:
            return True
        else:
            return False
    else:
        return False


def get_banner(value):
    banner = re.findall('(?im)<\s*banner.*>(.*?)<\s*/\s*banner>',rep_data)
    if banner == []:
        return False
    else:
        if value in str(banner):
            return True
        else:
            return False


def get_request(url):
    global rep_head
    global rep_data
    global req_head

    try:
        r = requests.get(url=url, headers=headers,verify=False)
        req_head = r.request.headers
        rep_head = r.headers
        r.encoding = "utf-8"
        rep_data = r.text
    except:
        print(url+"连接失败")

def get_title(value):
    title = re.findall('<title>(.*?)</title>',rep_data)
    if value in str(title):
        return True
    else:
        return False

def get_head(value):
    if value in str(rep_head) or value in str(req_head):
        return True
    else:
        return False

def get_server(value):
    if value in rep_head['Server']:
        return True
    else:
        return False

def get_body(value):
    if value in str(rep_data):
        return True
    else:
        return False


def get_cert(url,data):
    if isIP(url):
        return False
    else:
        try:
            hostname = urlparse(url)
            hostname = hostname.hostname
            print("检测证书中，较慢——————")
            c = ssl.create_default_context()
            s = c.wrap_socket(socket.socket(), server_hostname=hostname)
            s.connect((hostname, 443))
            cert = s.getpeercert()
            if data in str(cert):
                return True
            else:
                return False
        except:
            return False

if __name__ == '__main__':
    main()