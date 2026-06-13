// Dashboard API client
class DashboardAPI {
    async getTrend(platform, type, days) {
        const response = await fetch(`/api/trend?platform=${platform}&type=${type}&days=${days}`);
        return response.json();
    }

    async getPlatforms() {
        const response = await fetch('/api/platforms');
        return response.json();
    }

    async getSchedule() {
        const response = await fetch('/api/schedule');
        return response.json();
    }

    async getHistory(limit) {
        const response = await fetch(`/api/history?limit=${limit}`);
        return response.json();
    }
}

// Dashboard UI controller
class Dashboard {
    constructor() {
        this.api = new DashboardAPI();
        this.chart = null;
        this.init();
    }

    async init() {
        this.setupEventListeners();
        await this.loadAll();
        this.startAutoRefresh();
    }

    setupEventListeners() {
        document.getElementById('refresh-trend').addEventListener('click', () => this.loadTrend());
        document.getElementById('refresh-history').addEventListener('click', () => this.loadHistory());
    }

    async loadAll() {
        await Promise.all([
            this.loadPlatforms(),
            this.loadTrend(),
            this.loadSchedule(),
            this.loadHistory()
        ]);
    }

    async loadPlatforms() {
        try {
            const data = await this.api.getPlatforms();
            this.renderPlatforms(data.platforms);
        } catch (error) {
            console.error('Failed to load platforms:', error);
            document.getElementById('platforms-container').innerHTML =
                '<div class="error">加载失败</div>';
        }
    }

    renderPlatforms(platforms) {
        const container = document.getElementById('platforms-container');
        if (platforms.length === 0) {
            container.innerHTML = '<div class="empty">暂无平台</div>';
            return;
        }

        container.innerHTML = platforms.map(p => `
            <div class="card platform-card">
                <div class="platform-name">${this.getPlatformDisplayName(p.name)}</div>
                <div class="platform-status ${p.status}">${this.getStatusText(p.status)}</div>
            </div>
        `).join('');
    }

    getPlatformDisplayName(name) {
        const names = {
            'bilibili': 'Bilibili',
            'douyin': '抖音',
            'xiaohongshu': '小红书'
        };
        return names[name] || name;
    }

    getStatusText(status) {
        const statuses = {
            'active': '✓ 活跃',
            'inactive': '○ 离线'
        };
        return statuses[status] || status;
    }

    async loadTrend() {
        const platform = document.getElementById('platform-select').value;
        const type = document.getElementById('type-select').value;
        const days = document.getElementById('days-select').value;

        try {
            const data = await this.api.getTrend(platform, type, days);
            this.renderTrendChart(data);
            this.renderTrendStats(data);
        } catch (error) {
            console.error('Failed to load trend:', error);
        }
    }

    renderTrendChart(data) {
        const ctx = document.getElementById('trend-chart').getContext('2d');

        if (this.chart) {
            this.chart.destroy();
        }

        const labels = data.points.map(p => p.date);
        const values = data.points.map(p => p.value);

        this.chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: '数据量',
                    data: values,
                    borderColor: 'rgb(99, 102, 241)',
                    backgroundColor: 'rgba(99, 102, 241, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            stepSize: 1
                        }
                    }
                }
            }
        });
    }

    renderTrendStats(data) {
        const container = document.getElementById('trend-stats');
        const changeClass = data.change >= 0 ? 'positive' : 'negative';
        const changeSymbol = data.change >= 0 ? '↑' : '↓';

        container.innerHTML = `
            <div class="stat-card">
                <div class="stat-label">最小值</div>
                <div class="stat-value">${data.min}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">最大值</div>
                <div class="stat-value">${data.max}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">平均值</div>
                <div class="stat-value">${data.average.toFixed(2)}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">变化率</div>
                <div class="stat-value ${changeClass}">${changeSymbol} ${Math.abs(data.change).toFixed(2)}%</div>
            </div>
        `;
    }

    async loadSchedule() {
        try {
            const data = await this.api.getSchedule();
            this.renderSchedule(data.jobs);
        } catch (error) {
            console.error('Failed to load schedule:', error);
            document.getElementById('schedule-container').innerHTML =
                '<div class="error">加载失败</div>';
        }
    }

    renderSchedule(jobs) {
        const container = document.getElementById('schedule-container');
        if (jobs.length === 0) {
            container.innerHTML = '<div class="empty">暂无调度任务</div>';
            return;
        }

        const table = `
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>间隔</th>
                        <th>命令</th>
                        <th>上次运行</th>
                        <th>下次运行</th>
                        <th>状态</th>
                    </tr>
                </thead>
                <tbody>
                    ${jobs.map(job => `
                        <tr>
                            <td>${job.id}</td>
                            <td>${job.interval}</td>
                            <td class="command">${job.command}</td>
                            <td>${this.formatTime(job.last_run)}</td>
                            <td>${this.formatTime(job.next_run)}</td>
                            <td><span class="status-badge ${job.running ? 'running' : 'idle'}">${job.running ? '运行中' : '空闲'}</span></td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = table;
    }

    async loadHistory() {
        const limit = document.getElementById('history-limit').value;

        try {
            const data = await this.api.getHistory(limit);
            this.renderHistory(data.items);
        } catch (error) {
            console.error('Failed to load history:', error);
            document.getElementById('history-container').innerHTML =
                '<div class="error">加载失败</div>';
        }
    }

    renderHistory(items) {
        const container = document.getElementById('history-container');
        if (items.length === 0) {
            container.innerHTML = '<div class="empty">暂无历史记录</div>';
            return;
        }

        const table = `
            <table class="data-table">
                <thead>
                    <tr>
                        <th>平台</th>
                        <th>类型</th>
                        <th>采集时间</th>
                    </tr>
                </thead>
                <tbody>
                    ${items.map(item => `
                        <tr>
                            <td>${this.getPlatformDisplayName(item.platform)}</td>
                            <td><span class="type-badge">${item.type}</span></td>
                            <td>${this.formatTime(item.collected_at)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = table;
    }

    formatTime(timeStr) {
        if (!timeStr || timeStr === '0001-01-01T00:00:00Z') {
            return '-';
        }
        const date = new Date(timeStr);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    }

    startAutoRefresh() {
        setInterval(() => this.loadAll(), 30000); // 30秒自动刷新
    }
}

// Initialize dashboard when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new Dashboard();
});
