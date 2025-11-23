# 🧪 手动测试指南

由于 Clash-Fish 需要 **sudo 权限**来创建 TUN 设备和配置路由，请按以下步骤手动测试。

---

## 🚀 快速测试（推荐）

### 选项 A: 使用测试脚本（最简单）

```bash
# 运行交互式测试脚本
./scripts/test-manual.sh
```

脚本会引导你完成：
1. ✅ 配置验证
2. ✅ 启动前状态检查
3. ✅ 启动服务（需要输入 sudo 密码）
4. ✅ 按 Ctrl+C 停止服务
5. ✅ 验证清理完成

### 选项 B: 手动逐步测试

#### 第 1 步: 启动服务

```bash
sudo ./build/clash-fish start
```

**预期输出**：
```
[时间] INF Starting clash-fish service...
[时间] INF VPN detected, proxy will coexist with VPN interface=utun4 ip=10.8.0.105 network=10.8.0.0/24
✓ Clash-Fish started successfully
  Config: /Users/lg/.config/clash-fish/config.yaml

Service is running in foreground. Press Ctrl+C to stop.
```

#### 第 2 步: 在另一个终端检查状态

**打开新终端，运行检查脚本**：
```bash
cd /Users/lg/Projects/go/clash-fish
./scripts/check-running.sh
```

**或者手动检查**：
```bash
# 1. 检查服务状态
./build/clash-fish status

# 2. 检查 TUN 设备
ifconfig | grep utun5

# 3. 检查路由表
netstat -nr | grep 198.18

# 4. 检查端口
lsof -i :7890
lsof -i :7891

# 5. 测试代理
curl -x http://127.0.0.1:7890 -I https://www.google.com
```

#### 第 3 步: 测试 VPN 共存

```bash
# 查看路由表，验证 VPN 路由优先级
netstat -nr | head -30

# 应该看到：
# 10.8.0.0/24 -> utun4 (VPN，高优先级)
# 0.0.0.0/1   -> utun5 (代理，次优先级)
# default     -> en0   (默认，低优先级)
```

#### 第 4 步: 停止服务

**方式 1: 在运行终端**
```bash
按 Ctrl+C
```

**方式 2: 在另一个终端**
```bash
sudo ./build/clash-fish stop
```

**预期输出**：
```
[时间] INF Stopping clash-fish service...
✓ Clash-Fish stopped successfully
[时间] INF Service stopped pid=xxxxx
```

#### 第 5 步: 验证清理

```bash
# 检查状态
./build/clash-fish status
# 应该显示: Service: ✗ Not Running

# 检查 PID 文件
ls ~/.config/clash-fish/*.pid
# 应该显示: No such file or directory
```

---

## ✅ 测试检查清单

复制以下清单进行测试：

```
□ 配置验证通过
□ VPN 检测正确 (utun4: 10.8.0.105)
□ 服务成功启动
□ PID 文件创建 (~/.config/clash-fish/clash-fish.pid)
□ utun5 设备创建 (ifconfig)
□ 路由表正确配置 (0.0.0.0/1, 128.0.0.0/1 -> utun5)
□ VPN 路由优先级正确 (10.8.0.0/24 -> utun4)
□ HTTP 代理端口监听 (7890)
□ SOCKS5 代理端口监听 (7891)
□ 代理连接测试成功
□ 服务停止成功
□ PID 文件清理完成
```

---

## 🐛 常见问题

### Q1: 启动时提示 "requires root privileges"

**解决**: 必须使用 `sudo` 运行：
```bash
sudo ./build/clash-fish start
```

### Q2: 启动后没有任何输出

**检查**:
```bash
# 查看日志
tail -f ~/.config/clash-fish/logs/clash-fish.log

# 或者使用 debug 模式
sudo ./build/clash-fish start --debug
```

### Q3: 端口被占用

**解决**:
```bash
# 查找占用端口的进程
lsof -i :7890
lsof -i :7891

# 停止冲突的进程或修改配置文件端口
```

### Q4: VPN 流量被代理劫持

**检查路由表**:
```bash
netstat -nr | grep "10.8"
```

**应该看到**: `10.8.0.0/24 -> utun4`（优先级高于代理路由）

### Q5: 无法访问外网

**排查步骤**:
```bash
# 1. 检查 TUN 设备
ifconfig utun5

# 2. 检查路由
netstat -nr | grep 198.18

# 3. 检查代理服务器配置
cat ~/.config/clash-fish/config.yaml | grep -A5 "proxies:"

# 4. 测试代理连接
curl -x http://127.0.0.1:7890 -v https://www.google.com
```

---

## 📊 测试结果模板

测试完成后，可以记录结果：

```
测试日期: _______________
macOS 版本: _______________
VPN 状态: _______________

配置管理:
  □ config init:     PASS / FAIL
  □ config validate: PASS / FAIL
  □ config show:     PASS / FAIL

服务管理:
  □ start:   PASS / FAIL
  □ stop:    PASS / FAIL
  □ restart: PASS / FAIL
  □ status:  PASS / FAIL

网络功能:
  □ TUN 设备:    PASS / FAIL
  □ 路由配置:    PASS / FAIL
  □ HTTP 代理:   PASS / FAIL
  □ SOCKS5 代理: PASS / FAIL

VPN 共存:
  □ VPN 检测:        PASS / FAIL
  □ VPN 路由优先级:  PASS / FAIL
  □ 同时访问内外网:  PASS / FAIL

问题记录:
_______________________________________________
_______________________________________________
_______________________________________________
```

---

## 📖 相关文档

- [技术设计文档](./TECHNICAL_DESIGN.md) - 完整的技术设计
- [开发计划](./DEVELOPMENT_PLAN.md) - 开发计划和进度
- [测试报告](./TEST_REPORT.md) - 详细的测试报告

---

**准备好了吗？开始测试吧！** 🚀

如果遇到问题，请查看 [TEST_REPORT.md](./TEST_REPORT.md) 的故障排查部分。
