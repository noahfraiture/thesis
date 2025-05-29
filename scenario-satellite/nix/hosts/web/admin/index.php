<?php
// Generated with Grok
$uploadDir = 'uploads/';

if (!file_exists($uploadDir)) {
    mkdir($uploadDir, 0777, true);
}

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    $uploadFile = $uploadDir . basename($_FILES['userfile']['name']);
    
    if (move_uploaded_file($_FILES['userfile']['tmp_name'], $uploadFile)) {
        $uploadSuccess = "Firmware update package " . basename($_FILES['userfile']['name']) . " has been verified and queued for deployment.";
    } else {
        $uploadError = "Error: Failed to process the firmware package.";
    }
}
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OrionSatellite Systems - Mission Control</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { background-color: #0a192f; color: #ccd6f6; }
        .nav-gradient { background: linear-gradient(90deg, #0a192f 0%, #1a365f 100%); }
        .card-custom { background: #112240; border: 1px solid #233554; border-radius: 8px; }
        .chart-container { height: 300px; }
        .btn-space { background: #64ffda; color: #0a192f; }
        .stat-card { background: #233554; border-left: 4px solid #64ffda; }
    </style>
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark nav-gradient">
        <div class="container-fluid">
            <a class="navbar-brand" href="#">
                <i class="fas fa-satellite me-2"></i>
                OrionSatellite Systems
            </a>
        </div>
    </nav>

    <div class="container-fluid mt-4">
        <div class="row">
            <div class="col-md-3">
                <div class="card card-custom p-3 mb-4">
                    <h5 class="text-muted mb-3"><i class="fas fa-rocket me-2"></i>System Status</h5>
                    <div class="stat-card p-3 mb-3">
                        <div class="text-small">Online Satellites</div>
                        <div class="h4">14/16</div>
                    </div>
                    <div class="stat-card p-3 mb-3">
                        <div class="text-small">Signal Strength</div>
                        <div class="h4">98.7%</div>
                    </div>
                    <div class="stat-card p-3">
                        <div class="text-small">Orbit Stability</div>
                        <div class="h4">Optimal</div>
                    </div>
                </div>
            </div>

            <div class="col-md-9">
                <div class="card card-custom p-4 mb-4">
                    <h4 class="mb-4"><i class="fas fa-chart-line me-2"></i>Telemetry Overview</h4>
                    <div class="chart-container">
                        <canvas id="telemetryChart"></canvas>
                    </div>
                </div>

                <div class="card card-custom p-4 mb-4">
                    <h4 class="mb-4"><i class="fas fa-upload me-2"></i>Firmware Update Portal</h4>
                    <?php if(isset($uploadSuccess)): ?>
                        <div class="alert alert-success"><?= $uploadSuccess ?></div>
                    <?php endif; ?>
                    <?php if(isset($uploadError)): ?>
                        <div class="alert alert-danger"><?= $uploadError ?></div>
                    <?php endif; ?>
                    
                    <form action="" method="post" enctype="multipart/form-data">
                        <div class="mb-3">
                            <label class="form-label">Select Firmware Package</label>
                            <input type="file" class="form-control bg-dark text-light" name="userfile" required>
                        </div>
                        <button type="submit" class="btn btn-space">
                            <i class="fas fa-upload me-2"></i>Deploy Package
                        </button>
                    </form>
                </div>

                <div class="card card-custom p-4">
                    <h4 class="mb-4"><i class="fas fa-history me-2"></i>Upload Logs</h4>
                    <div class="table-responsive">
                        <table class="table table-dark">
                            <thead>
                                <tr>
                                    <th>Filename</th>
                                    <th>Upload Date</th>
                                    <th>Size</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                <?php
                                $files = scandir($uploadDir);
                                foreach($files as $file) {
                                    if($file === '.' || $file === '..') continue;
                                    $filePath = $uploadDir . $file;
                                    echo "<tr>
                                            <td>$file</td>
                                            <td>" . date('Y-m-d H:i', filemtime($filePath)) . "</td>
                                            <td>" . round(filesize($filePath)/1024, 2) . " KB</td>
                                            <td><a href='uploads/$file' class='btn btn-sm btn-space'>Inspect</a></td>
                                        </tr>";
                                }
                                ?>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Telemetry Chart
        const ctx = document.getElementById('telemetryChart').getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
                datasets: [{
                    label: 'Signal Quality (%)',
                    data: [92, 95, 97, 96, 98, 99],
                    borderColor: '#64ffda',
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { labels: { color: '#ccd6f6' } }
                },
                scales: {
                    x: { ticks: { color: '#8892b0' } },
                    y: { ticks: { color: '#8892b0' } }
                }
            }
        });
    </script>
</body>
</html>
