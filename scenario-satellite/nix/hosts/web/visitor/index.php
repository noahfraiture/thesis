<?php
// Generated with Grok
// Simple PHP header manipulation
header('Content-Type: text/html; charset=utf-8');
$company_name = "Stellar Horizon Innovations";
$tagline = "Pioneering the Future of Space Technology and Exploration";
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title><?= $company_name ?> - Space Technology & Tourism</title>
    <link href="https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&family=Poppins:wght@300;500;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --space-blue: #0B0D21;
            --star-white: #E6E6E6;
            --neon-blue: #00F3FF;
            --neon-purple: #BD00FF;
            --gradient: linear-gradient(45deg, var(--neon-blue), var(--neon-purple));
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Poppins', sans-serif;
            background: var(--space-blue);
            color: var(--star-white);
            overflow-x: hidden;
        }

        .stars {
            position: fixed;
            width: 100vw;
            height: 100vh;
            z-index: -1;
        }

        .star {
            position: absolute;
            background: white;
            border-radius: 50%;
            animation: twinkle var(--duration) ease-in-out infinite;
        }

        @keyframes twinkle {
            0%, 100% { opacity: 0.3; }
            50% { opacity: 1; }
        }

        .parallax {
            transform: translateZ(0);
            will-change: transform;
        }

        header {
            padding: 2rem;
            text-align: center;
            position: relative;
            overflow: hidden;
        }

        .hero {
            height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-direction: column;
            text-align: center;
            padding: 2rem;
            background: radial-gradient(circle at center, rgba(11, 13, 33, 0.8), rgba(11, 13, 33, 1));
        }

        h1 {
            font-family: 'Orbitron', sans-serif;
            font-size: 4rem;
            background: var(--gradient);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 2rem;
            text-shadow: 0 0 20px rgba(189, 0, 255, 0.3);
        }

        .cta-button {
            padding: 1rem 2rem;
            font-size: 1.2rem;
            background: var(--gradient);
            border: none;
            border-radius: 50px;
            color: var(--star-white);
            cursor: pointer;
            transition: transform 0.3s, box-shadow 0.3s;
            text-transform: uppercase;
            font-weight: bold;
        }

        .cta-button:hover {
            transform: translateY(-3px);
            box-shadow: 0 0 30px rgba(0, 243, 255, 0.5);
        }

        .section {
            padding: 5rem 2rem;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-direction: column;
            border-bottom: 1px solid rgba(230, 230, 230, 0.1);
        }

        .card-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 2rem;
            width: 100%;
            max-width: 1200px;
            margin-top: 3rem;
        }

        .card {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 20px;
            padding: 2rem;
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.1);
            transition: transform 0.3s;
        }

        .card:hover {
            transform: translateY(-10px);
        }

        .service-card {
            position: relative;
            overflow: hidden;
        }

        .service-card::before {
            content: '';
            position: absolute;
            top: -50%;
            left: -50%;
            width: 200%;
            height: 200%;
            background: var(--gradient);
            animation: rotate 4s linear infinite;
            opacity: 0.1;
        }

        @keyframes rotate {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .scroll-reveal {
            opacity: 0;
            transform: translateY(50px);
            transition: all 1s ease;
        }

        .visible {
            opacity: 1;
            transform: translateY(0);
        }

        .team-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 2rem;
            width: 100%;
            max-width: 1200px;
            margin-top: 3rem;
        }

        .team-member {
            text-align: center;
        }

        .team-member img {
            width: 150px;
            height: 150px;
            border-radius: 50%;
            object-fit: cover;
            margin-bottom: 1rem;
            border: 3px solid var(--neon-blue);
        }

        .timeline {
            position: relative;
            max-width: 1200px;
            margin: 3rem auto;
        }

        .timeline::before {
            content: '';
            position: absolute;
            top: 0;
            bottom: 0;
            width: 4px;
            background: var(--gradient);
            left: 50%;
            transform: translateX(-50%);
        }

        .timeline-item {
            padding: 20px 40px;
            position: relative;
            width: 50%;
        }

        .timeline-item:nth-child(odd) {
            left: 0;
        }

        .timeline-item:nth-child(even) {
            left: 50%;
        }

        .timeline-item::after {
            content: '';
            position: absolute;
            width: 20px;
            height: 20px;
            background: var(--gradient);
            border-radius: 50%;
            top: 20px;
            right: -10px;
        }

        .timeline-item:nth-child(even)::after {
            left: -10px;
        }
    </style>
</head>
<body>
    <div class="stars" id="stars"></div>

    <header class="parallax">
        <nav>
            <h2><?= $company_name ?></h2>
        </nav>
    </header>

    <section class="hero parallax">
        <h1 class="scroll-reveal">Beyond The Horizon</h1>
        <p class="scroll-reveal" style="font-size: 1.5rem; max-width: 800px; margin-bottom: 3rem;">
            <?= $tagline ?>
        </p>
        <button class="cta-button scroll-reveal">Launch Your Journey</button>
    </section>

    <section class="section">
        <h2 class="scroll-reveal">Our Services</h2>
        <div class="card-grid">
            <div class="card service-card scroll-reveal">
                <h3>🛰 Satellite Solutions</h3>
                <p>We design and deploy cutting-edge satellites for Earth observation, communication, and deep space exploration. Our satellites are equipped with AI-driven analytics for real-time data processing.</p>
            </div>
            <div class="card service-card scroll-reveal">
                <h3>🚀 Space Tourism</h3>
                <p>Experience the thrill of space travel with our suborbital flights, zero-gravity adventures, and lunar orbit packages. Safety and comfort are our top priorities.</p>
            </div>
            <div class="card service-card scroll-reveal">
                <h3>🌌 Space Research</h3>
                <p>Collaborate with us on groundbreaking research projects. From asteroid mining to interstellar travel, we push the boundaries of human knowledge.</p>
            </div>
            <div class="card service-card scroll-reveal">
                <h3>🛠 Spacecraft Manufacturing</h3>
                <p>We build state-of-the-art spacecraft for governments, private companies, and research institutions. Custom designs tailored to your mission needs.</p>
            </div>
        </div>
    </section>

    <section class="section">
        <h2 class="scroll-reveal">Our Mission</h2>
        <p class="scroll-reveal" style="max-width: 800px; text-align: center; font-size: 1.2rem;">
            At <?= $company_name ?>, our mission is to make space accessible to everyone. We believe in a future where humanity thrives among the stars, and we're committed to building the technology and infrastructure to make that vision a reality.
        </p>
    </section>

    <section class="section">
        <h2 class="scroll-reveal">Our Team</h2>
        <div class="team-grid">
            <div class="team-member scroll-reveal">
                <img src="https://via.placeholder.com/150" alt="Team Member">
                <h3>Dr. Elena Vega</h3>
                <p>Chief Scientist</p>
            </div>
            <div class="team-member scroll-reveal">
                <img src="https://via.placeholder.com/150" alt="Team Member">
                <h3>John Carter</h3>
                <p>Lead Engineer</p>
            </div>
            <div class="team-member scroll-reveal">
                <img src="https://via.placeholder.com/150" alt="Team Member">
                <h3>Sarah Orion</h3>
                <p>Space Tourism Director</p>
            </div>
            <div class="team-member scroll-reveal">
                <img src="https://via.placeholder.com/150" alt="Team Member">
                <h3>Alex Nova</h3>
                <p>Mission Control Specialist</p>
            </div>
        </div>
    </section>

    <section class="section">
        <h2 class="scroll-reveal">Our Timeline</h2>
        <div class="timeline">
            <div class="timeline-item scroll-reveal">
                <h3>2021 - Founded</h3>
                <p><?= $company_name ?> was established with a vision to revolutionize space technology.</p>
            </div>
            <div class="timeline-item scroll-reveal">
                <h3>2022 - First Satellite Launch</h3>
                <p>Successfully launched our first Earth observation satellite, Horizon-1.</p>
            </div>
            <div class="timeline-item scroll-reveal">
                <h3>2023 - Space Tourism Program</h3>
                <p>Announced our space tourism initiative, offering suborbital flights.</p>
            </div>
            <div class="timeline-item scroll-reveal">
                <h3>2024 - Lunar Mission</h3>
                <p>Planned lunar orbit mission for research and tourism.</p>
            </div>
        </div>
    </section>

    <section class="section">
        <h2 class="scroll-reveal">Contact Us</h2>
        <p class="scroll-reveal" style="max-width: 800px; text-align: center; font-size: 1.2rem;">
            Ready to explore the cosmos with us? Reach out to our team for inquiries, partnerships, or to book your space adventure.
        </p>
        <button class="cta-button scroll-reveal">Get in Touch</button>
    </section>

    <script>
        // Create animated stars background
        function createStars() {
            const stars = document.getElementById('stars');
            for (let i = 0; i < 200; i++) {
                const star = document.createElement('div');
                star.className = 'star';
                star.style.left = Math.random() * 100 + '%';
                star.style.top = Math.random() * 100 + '%';
                star.style.width = Math.random() * 3 + 'px';
                star.style.height = star.style.width;
                star.style.setProperty('--duration', Math.random() * 3 + 2 + 's');
                stars.appendChild(star);
            }
        }

        // Scroll reveal animation
        function handleScroll() {
            const elements = document.querySelectorAll('.scroll-reveal');
            elements.forEach(element => {
                const elementTop = element.getBoundingClientRect().top;
                if (elementTop < window.innerHeight - 100) {
                    element.classList.add('visible');
                }
            });
        }

        window.addEventListener('load', createStars);
        window.addEventListener('scroll', handleScroll);
        window.addEventListener('load', handleScroll);
    </script>
</body>
</html>
