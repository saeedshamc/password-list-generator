# نصب و اجرا (Installation & Run)

## نصب وابستگی‌های سیستم (System Dependencies)

برای اجرای GUI، این کتابخانه‌ها باید نصب باشند:

```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev libxxf86vm-dev
```

## ساخت برنامه (Build)

```bash
cd /home/hell-boy/Desktop/pass
go mod download
go build -o passgen .
```

## اجرای برنامه (Run)

```bash
./passgen
```

یا برای حالت CLI:

```bash
./passgen -cli
```

## اگر خطای لینک داشتید (Linker Errors)

اگر خطای `cannot find -lXxf86vm` دیدید، کتابخانه‌های بالا را نصب کنید.

## تست بدون GUI (Test without GUI)

اگر نمی‌خواهید GUI را نصب کنید، می‌توانید فقط core را تست کنید:

```bash
go test ./core/...
```

