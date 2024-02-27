import Chart from "chart.js/auto";

export function createCharts() {
  const canvasLine = document.getElementById('chartjs-line');
  const canvasPie = document.getElementById('chartjs-pie');
  if (!canvasLine || !canvasPie) return;

  const data = [
    { month: 'January', count: 10 },
    { month: 'February', count: 20 },
    { month: 'March', count: 15 },
    { month: 'April', count: 25 },
    { month: 'May', count: 22 },
    { month: 'June', count: 30 },
    { month: 'July', count: 28 },
  ];

  try {
    new Chart(
      canvasLine as HTMLCanvasElement,
      {
        type: 'line',
        data: {
          labels: data.map(row => row.month),
          datasets: [
            {
              data: data.map(row => row.count),
              fill: false,
              tension: 0.2,
              pointRadius: 5
            }
          ]
        },
        options: {
          responsive: true,
          onResize: (chart) => chart.resize(),
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: false,
            }
          },
          scales: {
            x: {
              grid: {
                display: false,
              }
            },
            y: {
              ticks: {
                maxTicksLimit: 5,
              }
            }
          }
        }
      }
    );

    new Chart(
      canvasPie as HTMLCanvasElement,
      {
        type: 'pie',
        data: {
          datasets: [
            {
              label: 'Purchases By Month',
              data: data.map(row => row.count)
            }
          ]
        },
        options: {
          responsive: true,
          onResize: (chart) => chart.resize(),
        }
      }
    );
  } catch (_) { }
}
