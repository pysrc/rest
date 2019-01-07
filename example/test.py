import requests

# 正确性验证
print(requests.get("http://127.0.0.1:8080/api/xiaom/12345/index").text)

print(requests.post("http://127.0.0.1:8080/api/xiaom/12345/index").text)

print(requests.put("http://127.0.0.1:8080/api/xiaom/12345/index").text)

print(requests.delete("http://127.0.0.1:8080/api/xiaom/12345/index").text)

# 不满足过滤规则验证

print(requests.get("http://127.0.0.1:8080/xapi/xiaom/12345/index").text)

# 不匹配验证
print(requests.post("http://127.0.0.1:8080/api/xiaom/index").text)

print(requests.put("http://127.0.0.1:8080/api/xiaom/12345").text)

print(requests.delete("http://127.0.0.1:8080/api/xiaom/12345/xxxx/index").text)

input()
