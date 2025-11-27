#================================= User Service =================================

# 1. Register
curl --location 'http://localhost:8080/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "example@mail.com",
    "full_name": "Belajar Golang Microservice",
    "password": "rahasia"
}'

# 2. Verify Register
curl --location 'http://localhost:8080/api/v1/auth/register/verify' \
--header 'Content-Type: application/json' \
--data '{
    "otp": "297135"
}'

# 3. Login
curl --location 'http://localhost:8080/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "example@mail.com",
    "password": "rahasia"
}'

# 4. Logout
curl --location --request DELETE 'http://localhost:8080/api/v1/auth/logout' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE'

# 5. Refresh Token
curl --location --request PATCH 'http://localhost:8080/api/v1/auth/refresh-token' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE'

# 6. Get Profile
curl --location 'http://localhost:8080/api/v1/profile' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE'



#================================= Product Service =================================

# 1. Create Product
curl --location 'http://localhost:8081/api/v1/products' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE' \
--form 'name="Semangka"' \
--form 'sku="Semangka1"' \
--form 'product_image=@"postman-cloud:///1eed8b0d-b243-4150-ab92-de33149a80a2"' \
--form 'price="5000"' \
--form 'stock="100"' \
--form 'length="15"' \
--form 'width="10"' \
--form 'height="5"' \
--form 'weight="0.5"' \
--form 'description="Semangka isi 3"'

# 2. Find Many Product
curl --location 'http://localhost:8081/api/v1/products' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE' \
--form 'name="Semangka"' \
--form 'sku="Semangka1"' \
--form 'product_image=@"postman-cloud:///1eed8b0d-b243-4150-ab92-de33149a80a2"' \
--form 'price="5000"' \
--form 'stock="100"' \
--form 'length="15"' \
--form 'width="10"' \
--form 'height="5"' \
--form 'weight="0.5"' \
--form 'description="Semangka isi 3"'



#================================= Order Service =================================
curl --location 'http://localhost:8082/api/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=YOUR_ACCESS_TOKEN_HERE; refresh_token=YOUR_REFRESH_TOKEN_HERE' \
--data '{
    "order": {
        "gross_amount": 10000
    },
    "product_orders": [
        {
            "product_id": 1,
            "quantity": 2,
            "price": 10000
        }
    ]
}'