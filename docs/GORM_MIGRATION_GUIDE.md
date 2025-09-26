# GORM AutoMigration Guide

## 🚀 **Tối ưu hóa Migration với GORM**

Sau khi tham khảo [GORM Migration documentation](https://gorm.io/docs/migration.html), chúng ta đã implement một hệ thống migration tối ưu hơn sử dụng GORM AutoMigrate.

## 📊 **So sánh các approaches:**

| Approach | Ưu điểm | Nhược điểm | Phù hợp |
|----------|---------|------------|---------|
| **Custom Migrator** (cũ) | Full control, Oracle-specific | Phải update code mỗi entity mới | Legacy systems |
| **GORM AutoMigrate** (mới) | Zero-config, struct-driven | Ít control hơn | Modern applications |
| **golang-migrate** | Version control, rollback | Phải viết SQL manually | Enterprise systems |

## ✅ **Ưu điểm của GORM AutoMigrate:**

### **1. Zero Configuration**
```go
// Chỉ cần 1 dòng code!
db.AutoMigrate(&User{}, &Product{}, &Order{})
```

### **2. Automatic Schema Management**
- ✅ Tạo tables tự động
- ✅ Tạo missing foreign keys
- ✅ Tạo constraints tự động
- ✅ Tạo columns tự động
- ✅ Tạo indexes tự động
- ✅ **KHÔNG xóa unused columns** (bảo vệ data)

### **3. Database Independent**
- ✅ Unified API cho tất cả databases
- ✅ Oracle support đầy đủ
- ✅ Automatic SQL generation

### **4. BMSF_ Prefix Support**
```go
// Custom naming strategy
type BMSFNamingStrategy struct{}

func (ns *BMSFNamingStrategy) TableName(table string) string {
    return "BMSF_" + strings.ToUpper(table)
}
```

## 🏗️ **Cách sử dụng:**

### **1. Define Entity với GORM Tags**
```go
type User struct {
    BaseEntity
    Username  string     `gorm:"size:50;uniqueIndex;not null"`
    Email     string     `gorm:"size:255;uniqueIndex;not null"`
    FirstName string     `gorm:"size:100;not null"`
    LastName  string     `gorm:"size:100;not null"`
    Phone     string     `gorm:"size:20"`
    Status    UserStatus `gorm:"size:20;default:'PENDING';not null"`
}
```

### **2. Run AutoMigration**
```go
// Tạo GORM migrator
migrator, err := database.NewGORMMigrator(dsn, logger)
if err != nil {
    log.Fatal(err)
}

// Auto-migrate tất cả entities
err = migrator.AutoMigrate(ctx)
if err != nil {
    log.Fatal(err)
}
```

### **3. Thêm Entity mới - KHÔNG CẦN UPDATE CODE!**
```go
// Chỉ cần thêm vào AutoMigrate call
err := m.db.AutoMigrate(
    &entities.User{},
    &entities.NewEntity{}, // ← Chỉ cần thêm dòng này!
)
```

## 📋 **GORM Tags Reference:**

### **Column Types:**
```go
`gorm:"type:varchar(50)"`           // VARCHAR2(50)
`gorm:"type:number(10,2)"`          // NUMBER(10,2)
`gorm:"type:timestamp"`             // TIMESTAMP
```

### **Constraints:**
```go
`gorm:"primaryKey"`                 // PRIMARY KEY
`gorm:"uniqueIndex"`                // UNIQUE INDEX
`gorm:"not null"`                   // NOT NULL
`gorm:"default:'value'"`            // DEFAULT 'value'
```

### **Indexes:**
```go
`gorm:"index"`                      // CREATE INDEX
`gorm:"uniqueIndex"`                // CREATE UNIQUE INDEX
`gorm:"index:idx_name"`             // Named index
```

### **Sizes:**
```go
`gorm:"size:50"`                    // VARCHAR2(50)
`gorm:"size:255"`                   // VARCHAR2(255)
```

## 🎯 **Kết quả:**

### **Tables được tạo:**
- `BMSF_USER` (với tất cả BaseEntity fields)

### **Indexes được tạo tự động:**
- `IDX_USER_USERNAME`
- `IDX_USER_EMAIL`
- `IDX_USER_DELETEDAT`
- `IDX_USER_TENANTID`
- Và nhiều indexes khác...

### **Constraints được tạo tự động:**
- `PK_USER_ID` (Primary Key)
- `UK_USER_USERNAME` (Unique)
- `UK_USER_EMAIL` (Unique)
- Và nhiều constraints khác...

## 🔧 **Migration Commands:**

### **Run Example:**
```bash
go run example_gorm_migration.go
```

### **Expected Output:**
```
✅ GORM AutoMigrate demo completed successfully!
📊 Created tables with BMSF_ prefix:
   - BMSF_USER (with all BaseEntity fields)
🔧 All indexes, constraints, and relationships created automatically!
🏷️  All table names follow BMSF_ prefix convention!
📋 All column names are in UPPERCASE (Oracle convention)!
```

## 🚀 **Lợi ích:**

1. **Zero Configuration**: Không cần viết SQL
2. **Automatic**: Tất cả được tạo tự động
3. **Safe**: Không xóa data
4. **Scalable**: Dễ dàng thêm entities mới
5. **Oracle Compatible**: Hỗ trợ đầy đủ Oracle
6. **BMSF_ Prefix**: Tuân thủ naming convention

## 📚 **References:**

- [GORM Migration Documentation](https://gorm.io/docs/migration.html)
- [GORM AutoMigrate](https://gorm.io/docs/migration.html#Auto-Migration)
- [GORM Naming Strategy](https://gorm.io/docs/conventions.html#Naming-Strategy)
