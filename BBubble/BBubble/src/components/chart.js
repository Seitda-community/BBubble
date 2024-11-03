    // 탭 표시 함수
    function showTab(tabId) {
        document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
        document.querySelectorAll('.content > div').forEach(content => content.classList.remove('active'));
        document.querySelector(`.tab[onclick="showTab('${tabId}')"]`).classList.add('active');
        document.getElementById(tabId).classList.add('active');
    }

    // Chart.js 그래프 생성 함수
    function renderCharts() {
        // 분야별 차트 데이터 및 설정
        new Chart(document.getElementById('fieldChart'), {
            type: 'bar',
            data: {
                labels: ['사이버', '대출', '고용', '임대', '폐업'],
                datasets: [{
                    label: '검색 횟수',
                    data: [150, 120, 90, 60, 30],
                    backgroundColor: '#523ae2'
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { display: false },
                    datalabels: {
                        anchor: 'end',
                        align: 'top',
                        color: '#000',
                        font: {
                            weight: 'bold'
                        },
                        formatter: (value) => `${value}회`
                    }
                }
            },
            plugins: [ChartDataLabels]
        });

        // 연령별 차트 데이터 및 설정
        new Chart(document.getElementById('ageChart'), {
            type: 'bubble',
            data: {
                datasets: [
                    { label: '전세사기', data: [{x: 20, y: 10, r: 45}], backgroundColor: '#523ae2' },
                    { label: '이혼', data: [{x: 30, y: 20, r: 30}], backgroundColor: '#db7195' },
                    { label: '대출', data: [{x: 40, y: 30, r: 22}], backgroundColor: '#75b2f2' },
                    { label: '임금', data: [{x: 50, y: 40, r: 18}], backgroundColor: '#ffc3a0' }
                ]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { display: false },
                    datalabels: {
                        color: '#fff',
                        font: {
                            weight: 'bold'
                        },
                        formatter: (value, context) => context.dataset.label
                    }
                },
                scales: {
                    x: { display: false },
                    y: { display: false }
                }
            },
            plugins: [ChartDataLabels]
        });

        // 지역별 차트 데이터 및 설정
        new Chart(document.getElementById('regionChart'), {
            type: 'doughnut',
            data: {
                labels: ['수원시', '성남시', '용인시', '안양시', '안산시', '과천시'],
                datasets: [{
                    data: [18.29, 17.18, 15.62, 12.28, 5.93, 5.79],
                    backgroundColor: ['#523ae2', '#db7195', '#a8d874', '#ffc3a0', '#5f5e5f', '#8f6ed9']
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { position: 'bottom' },
                    datalabels: {
                        color: '#fff',
                        font: {
                            weight: 'bold'
                        },
                        formatter: (value, context) => `${context.chart.data.labels[context.dataIndex]}: ${value}%`
                    }
                }
            },
            plugins: [ChartDataLabels]
        });
    }

    // 페이지 로드 후 차트 렌더링
    window.onload = function() {
        renderCharts();
    };