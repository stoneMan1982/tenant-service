// 基础配置
const API_BASE_URL = 'http://localhost:4300';

// DOM元素
const btnListMerchants = document.getElementById('btnListMerchants');
const btnCreateMerchant = document.getElementById('btnCreateMerchant');
const btnUpdateDomainPort = document.getElementById('btnUpdateDomainPort');
const btnSearchMerchant = document.getElementById('btnSearchMerchant');
const btnBackToList = document.getElementById('btnBackToList');

const merchantListPanel = document.getElementById('merchantList');
const createMerchantPanel = document.getElementById('createMerchant');
const updateDomainPortPanel = document.getElementById('updateDomainPort');
const merchantFilesPanel = document.getElementById('merchantFiles');

const merchantsTableBody = document.getElementById('merchantsTableBody');
const merchantFilesTableBody = document.getElementById('merchantFilesTableBody');
const currentMerchantId = document.getElementById('currentMerchantId');

const merchantListLoading = document.getElementById('merchantListLoading');
const merchantListError = document.getElementById('merchantListError');
const merchantFilesLoading = document.getElementById('merchantFilesLoading');
const merchantFilesError = document.getElementById('merchantFilesError');

const searchMerchantInput = document.getElementById('searchMerchant');
const createMerchantForm = document.getElementById('createMerchantForm');
const updateDomainPortForm = document.getElementById('updateDomainPortForm');
const createMerchantResult = document.getElementById('createMerchantResult');
const updateDomainPortResult = document.getElementById('updateDomainPortResult');

// 初始化
function init() {
    // 绑定事件处理程序
    btnListMerchants.addEventListener('click', () => switchPanel(merchantListPanel));
    btnCreateMerchant.addEventListener('click', () => switchPanel(createMerchantPanel));
    btnUpdateDomainPort.addEventListener('click', () => switchPanel(updateDomainPortPanel));
    btnSearchMerchant.addEventListener('click', searchMerchant);
    btnBackToList.addEventListener('click', () => switchPanel(merchantListPanel));
    
    createMerchantForm.addEventListener('submit', handleCreateMerchant);
    updateDomainPortForm.addEventListener('submit', handleUpdateDomainPort);
    
    // 搜索框回车事件
    searchMerchantInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            searchMerchant();
        }
    });
    
    // 初始加载商户列表
    loadMerchants();
}

// 切换面板
function switchPanel(panel) {
    // 隐藏所有面板
    merchantListPanel.classList.remove('active');
    createMerchantPanel.classList.remove('active');
    updateDomainPortPanel.classList.remove('active');
    merchantFilesPanel.classList.remove('active');
    
    // 移除所有导航按钮的激活状态
    btnListMerchants.classList.remove('active');
    btnCreateMerchant.classList.remove('active');
    btnUpdateDomainPort.classList.remove('active');
    
    // 显示选中的面板
    panel.classList.add('active');
    
    // 激活对应的导航按钮
    if (panel === merchantListPanel) {
        btnListMerchants.classList.add('active');
    } else if (panel === createMerchantPanel) {
        btnCreateMerchant.classList.add('active');
    } else if (panel === updateDomainPortPanel) {
        btnUpdateDomainPort.classList.add('active');
    }
    
    // 清除结果提示
    createMerchantResult.classList.add('hidden');
    updateDomainPortResult.classList.add('hidden');
}

// 加载商户列表
function loadMerchants() {
    merchantListLoading.style.display = 'block';
    merchantListError.classList.add('hidden');
    merchantsTableBody.innerHTML = '';
    
    fetch(`${API_BASE_URL}/merchants`)
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应错误');
            }
            return response.json();
        })
        .then(data => {
            merchantListLoading.style.display = 'none';
            
            if (data.success && data.data && data.data.length > 0) {
                data.data.forEach(merchant => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${merchant.merchantId}</td>
                        <td>
                            <button onclick="viewMerchantFiles('${merchant.merchantId}')" class="btn-secondary">查看文件</button>
                        </td>
                    `;
                    merchantsTableBody.appendChild(row);
                });
            } else {
                const row = document.createElement('tr');
                row.innerHTML = `<td colspan="2" style="text-align: center;">暂无商户数据</td>`;
                merchantsTableBody.appendChild(row);
            }
        })
        .catch(error => {
            merchantListLoading.style.display = 'none';
            merchantListError.classList.remove('hidden');
            merchantListError.textContent = '加载商户列表失败：' + error.message;
            console.error('加载商户列表失败:', error);
        });
}

// 搜索商户
function searchMerchant() {
    const searchTerm = searchMerchantInput.value.trim().toLowerCase();
    const rows = merchantsTableBody.querySelectorAll('tr');
    
    rows.forEach(row => {
        const merchantId = row.querySelector('td:first-child').textContent.toLowerCase();
        if (merchantId.includes(searchTerm)) {
            row.style.display = '';
        } else {
            row.style.display = 'none';
        }
    });
}

// 查看商户文件
function viewMerchantFiles(merchantId) {
    currentMerchantId.textContent = merchantId;
    switchPanel(merchantFilesPanel);
    
    merchantFilesLoading.style.display = 'block';
    merchantFilesError.classList.add('hidden');
    merchantFilesTableBody.innerHTML = '';
    
    fetch(`${API_BASE_URL}/merchant/${merchantId}/files`)
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应错误');
            }
            return response.json();
        })
        .then(data => {
            merchantFilesLoading.style.display = 'none';
            
            if (data.success && data.data && data.data.files && data.data.files.length > 0) {
                data.data.files.forEach(file => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${file.path}</td>
                        <td>${formatFileSize(file.size)}</td>
                        <td>${new Date(file.modTime || Date.now()).toLocaleString()}</td>
                        <td>
                            <button onclick="viewFileContent('${merchantId}', '${encodeURIComponent(file.path)}')" class="btn-secondary">查看内容</button>
                        </td>
                    `;
                    merchantFilesTableBody.appendChild(row);
                });
            } else {
                const row = document.createElement('tr');
                row.innerHTML = `<td colspan="4" style="text-align: center;">暂无文件数据</td>`;
                merchantFilesTableBody.appendChild(row);
            }
        })
        .catch(error => {
            merchantFilesLoading.style.display = 'none';
            merchantFilesError.classList.remove('hidden');
            merchantFilesError.textContent = '加载文件列表失败：' + error.message;
            console.error('加载文件列表失败:', error);
        });
}

// 查看文件内容
function viewFileContent(merchantId, filePath) {
    const decodedPath = decodeURIComponent(filePath);
    
    // 显示加载状态
    const fileContentModal = document.getElementById('fileContentModal');
    const fileContentLoading = document.getElementById('fileContentLoading');
    const fileContentError = document.getElementById('fileContentError');
    const fileContentTitle = document.getElementById('fileContentTitle');
    const fileContentBody = document.getElementById('fileContentBody');
    
    // 如果模态框不存在，创建它
    if (!fileContentModal) {
        createFileContentModal();
        return viewFileContent(merchantId, filePath);
    }
    
    fileContentTitle.textContent = `文件内容: ${decodedPath}`;
    fileContentBody.textContent = '';
    fileContentLoading.style.display = 'block';
    fileContentError.classList.add('hidden');
    
    // 显示模态框
    fileContentModal.style.display = 'block';
    
    // 发送请求获取文件内容
    fetch(`${API_BASE_URL}/merchant/${merchantId}/file?path=${encodeURIComponent(filePath)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`网络响应错误: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            fileContentLoading.style.display = 'none';
            
            if (data.success && data.data) {
                // 根据文件扩展名设置适当的内容格式
                const fileExtension = getFileExtension(decodedPath).toLowerCase();
                const content = data.data.content;
                
                if (['html', 'htm', 'css', 'js', 'json', 'xml', 'txt'].includes(fileExtension)) {
                    // 对于文本文件，显示格式化的内容
                    fileContentBody.innerHTML = `<pre><code>${escapeHtml(content)}</code></pre>`;
                } else {
                    // 对于二进制文件或不支持的格式，显示信息
                    fileContentBody.innerHTML = `<p>文件类型: ${fileExtension}</p>
                        <p>文件大小: ${formatFileSize(data.data.size)}</p>
                        <p>此文件类型可能包含二进制内容，无法直接显示</p>`;
                }
            } else {
                fileContentError.classList.remove('hidden');
                fileContentError.textContent = `获取文件内容失败: ${data.error || '未知错误'}`;
            }
        })
        .catch(error => {
            fileContentLoading.style.display = 'none';
            fileContentError.classList.remove('hidden');
            fileContentError.textContent = `获取文件内容失败: ${error.message}`;
            console.error('获取文件内容失败:', error);
        });
}

// 创建文件内容模态框
function createFileContentModal() {
    const modal = document.createElement('div');
    modal.id = 'fileContentModal';
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3 id="fileContentTitle">文件内容</h3>
                <button id="closeFileContentModal" class="close-btn">&times;</button>
            </div>
            <div class="modal-body">
                <div id="fileContentLoading" class="loading">加载中...</div>
                <div id="fileContentError" class="error hidden">获取文件内容失败</div>
                <div id="fileContentBody" class="file-content"></div>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    // 绑定关闭事件
    document.getElementById('closeFileContentModal').addEventListener('click', () => {
        document.getElementById('fileContentModal').style.display = 'none';
    });
    
    // 点击模态框外部关闭
    window.addEventListener('click', (event) => {
        const modal = document.getElementById('fileContentModal');
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}

// 获取文件扩展名
function getFileExtension(filename) {
    const parts = filename.split('.');
    if (parts.length > 1) {
        return parts[parts.length - 1];
    }
    return '';
}

// 转义HTML特殊字符
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 处理创建商户
function handleCreateMerchant(e) {
    e.preventDefault();
    
    const merchantId = document.getElementById('merchantId').value.trim();
    const domain = document.getElementById('domain').value.trim() || 'localhost';
    const port = document.getElementById('port').value.trim() || '8080';
    const protocol = document.getElementById('protocol').value;
    
    if (!merchantId) {
        showResult(createMerchantResult, false, '商户ID不能为空');
        return;
    }
    
    // 构建查询参数
    const params = new URLSearchParams();
    params.append('domain', domain);
    params.append('port', port);
    params.append('protocol', protocol);
    
    fetch(`${API_BASE_URL}/merchant/create/${merchantId}?${params.toString()}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应错误');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showResult(createMerchantResult, true, '商户创建成功');
                createMerchantForm.reset();
                // 刷新商户列表
                loadMerchants();
            } else {
                showResult(createMerchantResult, false, data.error || '商户创建失败');
            }
        })
        .catch(error => {
            showResult(createMerchantResult, false, '商户创建失败：' + error.message);
            console.error('创建商户失败:', error);
        });
}

// 处理更新域名和端口
function handleUpdateDomainPort(e) {
    e.preventDefault();
    
    const merchantId = document.getElementById('updateMerchantId').value.trim();
    const domain = document.getElementById('updateDomain').value.trim();
    const port = document.getElementById('updatePort').value.trim();
    
    if (!merchantId || !domain || !port) {
        showResult(updateDomainPortResult, false, '商户ID、域名和端口不能为空');
        return;
    }
    
    // 构建查询参数
    const params = new URLSearchParams();
    params.append('domain', domain);
    params.append('port', port);
    
    fetch(`${API_BASE_URL}/merchant/${merchantId}/domain-port?${params.toString()}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应错误');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showResult(updateDomainPortResult, true, `域名和端口更新成功，共处理 ${data.data.filesProcessed} 个文件`);
            } else {
                showResult(updateDomainPortResult, false, data.error || '更新失败');
            }
        })
        .catch(error => {
            showResult(updateDomainPortResult, false, '更新失败：' + error.message);
            console.error('更新域名和端口失败:', error);
        });
}

// 显示结果提示
function showResult(element, success, message) {
    element.textContent = message;
    element.className = 'result ' + (success ? 'success' : 'error');
    element.classList.remove('hidden');
    
    // 3秒后自动隐藏
    setTimeout(() => {
        element.classList.add('hidden');
    }, 3000);
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 当页面加载完成后初始化
document.addEventListener('DOMContentLoaded', init);