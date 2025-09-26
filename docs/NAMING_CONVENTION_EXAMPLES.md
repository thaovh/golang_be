# BMSF Naming Convention Examples

## 🏷️ **Updated Naming Strategy - Tránh trùng lặp giữa các tables**

### **Vấn đề trước đây:**
```
❌ IDX_USER_DELETEDAT    (từ BMSF_USER)
❌ IDX_PRODUCT_DELETEDAT (từ BMSF_PRODUCT) 
❌ IDX_ORDER_DELETEDAT   (từ BMSF_ORDER)
```
**→ Tất cả đều có column `DELETEDAT` từ BaseEntity → Trùng lặp!**

### **Giải pháp mới:**
```
✅ IDX_USER_DELETEDAT     (từ BMSF_USER table)
✅ IDX_PRODUCT_DELETEDAT  (từ BMSF_PRODUCT table)
✅ IDX_ORDER_DELETEDAT    (từ BMSF_ORDER table)
```
**→ Mỗi table có index name riêng biệt, không trùng lặp!**

## 📊 **Naming Convention Examples:**

### **Tables:**
```
BMSF_USER     → BMSF_USER
```

### **Indexes (tránh trùng lặp):**
```
BMSF_USER table:
- IDX_USER_ID          (Primary key index)
- IDX_USER_USERNAME    (Unique index)
- IDX_USER_EMAIL       (Unique index)
- IDX_USER_DELETEDAT   (Soft delete index)
- IDX_USER_TENANTID    (Multi-tenant index)
```

### **Constraints (tránh trùng lặp):**
```
BMSF_USER table:
- PK_USER_ID           (Primary key)
- UK_USER_USERNAME     (Unique constraint)
- UK_USER_EMAIL        (Unique constraint)
```

### **Foreign Keys (tránh trùng lặp):**
```
No foreign keys in current implementation
```

## 🔧 **Oracle 30-Character Limit Handling:**

### **Short Names (≤ 30 chars):**
```
✅ IDX_USER_USERNAME     (18 chars)
✅ UK_USER_EMAIL         (15 chars)
✅ FK_ORDER_USERID       (16 chars)
```

### **Long Names (> 30 chars) - Auto Truncate:**
```
❌ IDX_VERYLONGTABLENAME_VERYLONGCOLUMNNAME (45 chars)
✅ IDX_VERYLONGTABLENAME_VERYLO (30 chars) - Truncated

❌ UK_VERYLONGTABLENAME_VERYLONGCOLUMNNAME (45 chars)  
✅ UK_VERYLONGTABLENAME_VERYLO (30 chars) - Truncated
```

## 🎯 **Benefits:**

### **1. No Conflicts:**
- ✅ Mỗi table có index/constraint names riêng biệt
- ✅ Không trùng lặp giữa các tables
- ✅ Dễ dàng identify table source

### **2. Oracle Compatible:**
- ✅ Tuân thủ 30-char limit
- ✅ Auto truncate khi cần thiết
- ✅ Giữ table name để maintain uniqueness

### **3. Readable:**
- ✅ Clear naming pattern: `TYPE_TABLE_COLUMN`
- ✅ Easy to understand và debug
- ✅ Consistent across all objects

## 📝 **Naming Pattern:**

```
Format: {TYPE}_{TABLE}_{COLUMN}

Types:
- IDX_  → Index
- UK_   → Unique Constraint  
- FK_   → Foreign Key
- CHK_  → Check Constraint
- PK_   → Primary Key (handled by GORM)

Examples:
- IDX_USER_USERNAME    → Index on USER.USERNAME
- UK_PRODUCT_CODE      → Unique constraint on PRODUCT.CODE
- FK_ORDER_USERID      → Foreign key ORDER.USERID → USER.ID
```

## 🚀 **Migration Impact:**

### **Before (Conflicting):**
```sql
-- Multiple tables with same index names
CREATE INDEX IDX_DELETEDAT ON BMSF_USER(DELETEDAT);     -- ❌ Conflict
CREATE INDEX IDX_DELETEDAT ON BMSF_PRODUCT(DELETEDAT);  -- ❌ Conflict  
CREATE INDEX IDX_DELETEDAT ON BMSF_ORDER(DELETEDAT);    -- ❌ Conflict
```

### **After (Unique):**
```sql
-- Each table has unique index names
CREATE INDEX IDX_USER_DELETEDAT ON BMSF_USER(DELETEDAT);     -- ✅ Unique
CREATE INDEX IDX_PRODUCT_DELETEDAT ON BMSF_PRODUCT(DELETEDAT); -- ✅ Unique
CREATE INDEX IDX_ORDER_DELETEDAT ON BMSF_ORDER(DELETEDAT);   -- ✅ Unique
```

**→ No more naming conflicts!** 🎉
