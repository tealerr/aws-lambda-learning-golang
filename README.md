# Go Hello World on AWS Lambda

ตัวอย่างแอป Go + Gin ที่มี endpoint พื้นฐาน 2 ตัว และรองรับการรันทั้งบนเครื่องและ AWS Lambda

| Method | Path | ผลลัพธ์ |
| --- | --- | --- |
| `GET` | `/` | `{"message":"Hello, World!"}` |
| `GET` | `/health` | `{"status":"ok"}` |

## สิ่งที่ต้องติดตั้ง

- Go 1.25 ขึ้นไป
- AWS SAM CLI สำหรับรันหรือ deploy Lambda
- Docker สำหรับ `sam local`

## รันบนเครื่อง

```sh
go run main.go
```

เมื่อรันนอก Lambda แอปจะเปิดที่ `http://localhost:8080` และแสดง log พร้อม endpoint:

```text
🚀 Local server running at http://localhost:8080
📍 Endpoints: GET / (Hello World), GET /health (Health check)
```

ทดสอบด้วย:

```sh
curl http://localhost:8080/
curl http://localhost:8080/health
```

หรือเปิดโหมด local อย่างชัดเจนด้วย environment variable `LOCAL_SERVER=true` ก่อนรัน

## รัน Lambda ในเครื่องด้วย SAM

```sh
sam build
sam local start-api
```

เรียก API ที่ SAM เปิดไว้:

```sh
curl http://127.0.0.1:3000/Prod/
curl http://127.0.0.1:3000/Prod/health
```

## Deploy ไป AWS

```sh
sam build
sam deploy --guided \
  --parameter-overrides \
  'DatadogExtensionLayerArn=arn:aws:lambda:REGION:464622532012:layer:Datadog-Extension-ARM:91' \
  'DatadogApiKeySecretArn=SECRET_ARN' \
  DatadogSite=datadoghq.com
```

ก่อน deploy ให้สร้าง AWS Secrets Manager secret ที่เก็บ Datadog API key เป็น plaintext string และใช้ Extension layer ARN ที่ตรงกับ region และ architecture `arm64` ของฟังก์ชัน (ตัวอย่าง ARN ใช้ AWS commercial region และ layer version 91 ตาม Datadog docs) execution role จะได้รับ `secretsmanager:GetSecretValue` สำหรับ secret นี้จาก template

กำหนด `DatadogSite` ให้ตรงกับ Datadog site ของบัญชี เช่น `datadoghq.com`, `us3.datadoghq.com` หรือ `datadoghq.eu` เมื่อ deploy เสร็จ ใช้ API Gateway URL ที่ SAM แสดง แล้วเรียก `/` หรือ `/health`

ฟังก์ชันส่ง unified service tags ผ่าน `DD_ENV`, `DD_SERVICE` และ `DD_VERSION` ซึ่งกำหนดค่าเริ่มต้นไว้ใน `template.yaml` และสามารถแก้ให้ตรงกับ environment/release ของคุณได้

## โครงสร้างไฟล์

- `main.go` — endpoint และตัวเริ่ม Gin server หรือ Lambda runtime
- `template.yaml` — AWS SAM template สำหรับ Lambda + API Gateway (`provided.al2023`, `arm64`)
- `Dockerfile` — สร้าง Lambda container image
- `go.mod` และ `go.sum` — Go module และ dependencies
