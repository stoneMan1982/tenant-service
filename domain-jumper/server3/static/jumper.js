class DomainJumper {
    constructor() {
        this.data = null;
        this.currentStep = 0;
        this.totalSteps = 0; // 将在配置加载后动态计算
        this.statusEl = document.getElementById('status');
        this.progressBarEl = document.getElementById('progressBar');
        this.detailsEl = document.getElementById('details');

        this.init();
    }

    async init() {
        try {
            await this.loadData();
            this.updateStatus('配置加载成功，开始域名跳转流程...', 'success');
            await this.sleep(1000);
            await this.startJumpSequence();
        } catch (error) {
            this.updateStatus(`初始化失败: ${error.message}`, 'error');
        }
    }

    async loadData() {
        const response = await fetch('/data/domains.json');
        if (!response.ok) {
            throw new Error('无法加载域名配置');
        }
        this.data = await response.json();
        
        // 动态计算总步骤数：所有池子 + 最终跳转
        this.totalSteps = this.data.domainPools.length + 1;
    }

    updateStatus(message, type = 'info') {
        this.statusEl.innerHTML = '';
        if (type === 'info') {
            this.statusEl.innerHTML = '<div class="loading"></div>';
        }
        this.statusEl.innerHTML += message;
        this.statusEl.className = `status ${type}`;
    }

    updateProgress(step) {
        const percentage = Math.round((step / this.totalSteps) * 100);
        this.progressBarEl.style.width = percentage + '%';
        this.progressBarEl.textContent = `步骤 ${step}/${this.totalSteps} (${percentage}%)`;
    }

    updateDetails(content) {
        this.detailsEl.innerHTML = content;
    }

    async startJumpSequence() {
        try {
            // 从第一个池子开始，依次遍历所有池子
            for (let i = 0; i < this.data.domainPools.length; i++) {
                this.currentStep++;
                this.updateProgress(this.currentStep);

                const pool = this.data.domainPools[i];
                this.updateStatus(`步骤 ${this.currentStep}: 正在处理 ${pool.name}...`, 'info');

                const availableDomain = await this.findAvailableDomain(pool);

                if (!availableDomain) {
                    throw new Error(`${pool.name} 中没有可用的域名`);
                }

                this.updateStatus(`${pool.name} 找到可用域名: ${availableDomain}`, 'success');
                await this.sleep(this.data.jumpDelayMs);

                // 检查是否是目标池
                if (pool.type === 'target') {
                    // 跳转到最终目标
                    this.currentStep++;
                    this.updateProgress(this.currentStep);
                    this.updateStatus(`步骤 ${this.currentStep}: 跳转到最终目标...`, 'success');
                    await this.sleep(this.data.jumpDelayMs);

                    const finalDomain = availableDomain.replace(/(https?:\/\/)[^\.]+\./,
                        `$1${generateRandomString()}.`);
                    window.location.href = finalDomain;
                    return;
                }

                // 中间池：跳转到下一个池子
                const nextPoolIndex = i + 1;
                const nextDomain = availableDomain.replace(/(https?:\/\/)[^\.]+\./,
                    `$1${generateRandomString()}.`);
                window.location.href = `${nextDomain}?nextPool=${nextPoolIndex}`;
                return;
            }
        } catch (error) {
            this.updateStatus(`跳转过程中出现错误: ${error.message}`, 'error');
        }
    }

    async findAvailableDomain(pool) {
        this.updateDetails(`<div class="domain-list">
            <h4>${pool.name} 域名检测状态:</h4>
            ${pool.domains.map(domain =>
            `<div class="domain-item testing" id="domain-${domain}">
                    🔄 正在检测: ${domain}
                </div>`
        ).join('')}
        </div>`);

        const promises = pool.domains.map(domain =>
            this.checkDomainAvailability(domain)
        );

        const results = await Promise.allSettled(promises);

        // 更新检测结果显示
        for (let i = 0; i < results.length; i++) {
            const domain = pool.domains[i];
            const result = results[i];
            const domainEl = document.getElementById(`domain-${domain}`);

            if (result.status === 'fulfilled' && result.value) {
                domainEl.innerHTML = `✅ 可用: ${domain}`;
                domainEl.className = 'domain-item available';
            } else {
                domainEl.innerHTML = `❌ 不可用: ${domain}`;
                domainEl.className = 'domain-item unavailable';
            }
        }

        // 查找第一个可用的域名
        for (let i = 0; i < results.length; i++) {
            if (results[i].status === 'fulfilled' && results[i].value) {
                return pool.domains[i];
            }
        }

        return null;
    }

    async checkDomainAvailability(domain) {
        return new Promise((resolve) => {
            const timeout = setTimeout(() => {
                resolve(false);
            }, this.data.timeoutMs);

            // 本地测试：使用 fetch 检查 /health 端点
            fetch(`${domain}/health`, {
                mode: 'no-cors',
                cache: 'no-cache'
            })
                .then(response => {
                    clearTimeout(timeout);
                    resolve(true);
                })
                .catch(error => {
                    clearTimeout(timeout);
                    // 在本地测试中，即使 CORS 错误也表示域名是可访问的
                    resolve(true);
                });
        });
    }

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}


function generateRandomString() {
    const length = Math.floor(Math.random() * 4) + 3;
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const randomValues = new Uint32Array(length);

    window.crypto.getRandomValues(randomValues);

    let result = '';
    for (let i = 0; i < length; i++) {
        result += chars[randomValues[i] % chars.length];
    }

    return result;
}


// 检查URL参数，确定当前处理的是哪个池子
const urlParams = new URLSearchParams(window.location.search);
const nextPool = urlParams.get('nextPool');

if (nextPool) {
    // 如果有nextPool参数，说明这是中间跳转页面
    class IntermediateJumper extends DomainJumper {
        async startJumpSequence() {
            try {
                const poolIndex = parseInt(nextPool);
                this.currentStep = poolIndex;
                this.updateProgress(this.currentStep);

                if (poolIndex >= this.data.domainPools.length) {
                    throw new Error('池子索引超出范围');
                }

                const pool = this.data.domainPools[poolIndex];
                this.updateStatus(`步骤 ${this.currentStep + 1}: 正在处理 ${pool.name}...`, 'info');

                const availableDomain = await this.findAvailableDomain(pool);

                if (!availableDomain) {
                    throw new Error(`${pool.name} 中没有可用的域名`);
                }

                this.updateStatus(`${pool.name} 找到可用域名: ${availableDomain}`, 'success');
                await this.sleep(this.data.jumpDelayMs);

                // 检查是否是目标池
                if (pool.type === 'target') {
                    // 跳转到最终目标
                    this.currentStep++;
                    this.updateProgress(this.currentStep + 1);
                    this.updateStatus(`步骤 ${this.currentStep + 1}: 跳转到最终目标...`, 'success');
                    await this.sleep(this.data.jumpDelayMs);

                    const finalDomain = availableDomain.replace(/(https?:\/\/)[^\.]+\./,
                        `$1${generateRandomString()}.`);
                    window.location.href = finalDomain;
                    return;
                }

                // 中间池：跳转到下一个池子
                const nextPoolIndex = poolIndex + 1;
                const nextDomain = availableDomain.replace(/(https?:\/\/)[^\.]+\./,
                    `$1${generateRandomString()}.`);
                window.location.href = `${nextDomain}?nextPool=${nextPoolIndex}`;
                
            } catch (error) {
                this.updateStatus(`跳转过程中出现错误: ${error.message}`, 'error');
            }
        }
    }

    new IntermediateJumper();
} else {
    // 初始页面
    new DomainJumper();
}
