document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('recipe-form');
    const submitBtn = document.getElementById('submit-btn');
    const inputSection = document.getElementById('input-section');
    const resultSection = document.getElementById('result-section');
    const backBtn = document.getElementById('back-btn');

    // Timer Elements
    const timerToggleBtn = document.getElementById('timer-toggle-btn');
    const timerResetBtn = document.getElementById('timer-reset-btn');
    const timerMinutes = document.getElementById('timer-minutes');
    const timerSeconds = document.getElementById('timer-seconds');

    // Phase 3 Elements
    const flavorAcidity = document.getElementById('flavor-acidity');
    const flavorSweetness = document.getElementById('flavor-sweetness');
    const flavorBody = document.getElementById('flavor-body');
    const flavorBitterness = document.getElementById('flavor-bitterness');

    const copyRecipeBtn = document.getElementById('copy-recipe-btn');

    const openCalcBtn = document.getElementById('open-calc-btn');
    const closeCalcBtn = document.getElementById('close-calc-btn');
    const calcModal = document.getElementById('calc-modal');
    const applyCalcBtn = document.getElementById('apply-calc-btn');
    const calcTargetRatio = document.getElementById('calc-target-ratio');
    const calcCoffeeDose = document.getElementById('calc-coffee-dose');
    const calcWaterResult = document.getElementById('calc-water-result');

    // History Elements
    const historySection = document.getElementById('history-section');
    const historyContainer = document.getElementById('history-container');
    const openHistoryBtn = document.getElementById('open-history-btn');
    const historyBackBtn = document.getElementById('history-back-btn');
    const clearHistoryBtn = document.getElementById('clear-history-btn');

    // DOM Elements for Result
    const totalScoreText = document.getElementById('total-score-text');
    const scoreCirclePath = document.getElementById('score-circle-path');

    const calcRatioSpan = document.getElementById('calc-ratio');

    const idealRatioSpan = document.getElementById('ideal-ratio');
    const idealTempSpan = document.getElementById('ideal-temp');
    const idealTimeSpan = document.getElementById('ideal-time');

    const ratioProgress = document.getElementById('ratio-progress');
    const tempProgress = document.getElementById('temp-progress');
    const timeProgress = document.getElementById('time-progress');

    const feedbackList = document.getElementById('feedback-list');

    // Default values mapping
    const defaults = {
        'v60': { dose: 15, yield: 250, temp: 93, time: 180, grind: 'medium-fine', roast: 'light' },
        'espresso': { dose: 18, yield: 36, temp: 92, time: 30, grind: 'fine', roast: 'medium' },
        'aeropress': { dose: 15, yield: 200, temp: 85, time: 120, grind: 'medium', roast: 'light' },
        'frenchpress': { dose: 20, yield: 300, temp: 94, time: 270, grind: 'coarse', roast: 'medium' },
        'coldbrew': { dose: 50, yield: 500, temp: 20, time: 57600, grind: 'coarse', roast: 'dark' },
    };

    document.getElementById('method').addEventListener('change', (e) => {
        const method = e.target.value;
        const vals = defaults[method];
        if (vals) {
            document.getElementById('coffee_dose').value = vals.dose;
            document.getElementById('water_yield').value = vals.yield;
            document.getElementById('temperature').value = vals.temp;
            document.getElementById('brew_time').value = vals.time;
            document.getElementById('grind_size').value = vals.grind;
            document.getElementById('roast_level').value = vals.roast;
        }
    });

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        // UI feedback
        submitBtn.classList.add('loading');

        const formData = new FormData(form);
        const data = {
            method: formData.get('method'),
            coffee_dose: parseFloat(formData.get('coffee_dose')),
            water_yield: parseFloat(formData.get('water_yield')),
            temperature: parseFloat(formData.get('temperature')),
            brew_time: parseInt(formData.get('brew_time')),
            grind_size: formData.get('grind_size'),
            roast_level: formData.get('roast_level')
        };

        try {
            const response = await fetch('/api/score', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) throw new Error('Network response was not ok');

            const result = await response.json();

            // Artificial delay for UI polish
            setTimeout(() => {
                submitBtn.classList.remove('loading');
                showResult(result);
                saveToHistory(data, result);
            }, 500);

        } catch (error) {
            console.error('Error:', error);
            submitBtn.classList.remove('loading');
            alert('Terjadi kesalahan saat mengevaluasi skor. Pastikan server nyala.');
        }
    });

    backBtn.addEventListener('click', () => {
        resultSection.classList.add('hidden');
        inputSection.classList.remove('hidden');

        // Reset animations
        scoreCirclePath.style.strokeDasharray = `0, 100`;
        totalScoreText.textContent = '0';
        ratioProgress.style.width = '0%';
        tempProgress.style.width = '0%';
        timeProgress.style.width = '0%';
    });

    function showResult(result) {
        inputSection.classList.add('hidden');
        resultSection.classList.remove('hidden');

        // Populate Ideals
        idealRatioSpan.textContent = result.ideal_ratio;
        idealTempSpan.textContent = result.ideal_temp;
        idealTimeSpan.textContent = result.ideal_time;

        calcRatioSpan.textContent = `1:${result.calculated_ratio.toFixed(1)}`;

        // Populate Feedback (staggered animation)
        feedbackList.innerHTML = '';
        result.feedback.forEach((text, index) => {
            const li = document.createElement('li');
            li.textContent = text;
            li.style.animationDelay = `${(0.5 + (index * 0.1)).toFixed(1)}s`;
            feedbackList.appendChild(li);
        });

        // Trigger Animations (small delay to allow DOM render)
        setTimeout(() => {
            animateScore(result.total_score);

            // Progress bars (calculate percentage based on max score per category)
            const ratioPct = (result.ratio_score / 40) * 100;
            const tempPct = (result.temp_score / 30) * 100;
            const timePct = (result.time_score / 30) * 100;

            ratioProgress.style.width = `${ratioPct}%`;
            tempProgress.style.width = `${tempPct}%`;
            timeProgress.style.width = `${timePct}%`;

            // Adjust circle dash array (max 100)
            scoreCirclePath.style.strokeDasharray = `${result.total_score}, 100`;

            // Color mapping based on score
            if (result.total_score >= 90) {
                scoreCirclePath.parentNode.style.stroke = 'var(--accent-gold)';
            } else if (result.total_score >= 80) {
                scoreCirclePath.parentNode.style.stroke = '#a0d468';
            } else {
                scoreCirclePath.parentNode.style.stroke = 'var(--error)';
            }

        }, 50);
    }

    function animateScore(targetScore) {
        let currentScore = 0;
        const duration = 1500; // ms
        const steps = 60;
        const stepTime = Math.abs(Math.floor(duration / steps));
        const increment = targetScore / steps;

        const timer = setInterval(() => {
            currentScore += increment;
            if (currentScore >= targetScore) {
                currentScore = targetScore;
                clearInterval(timer);
            }
            totalScoreText.textContent = Math.round(currentScore);
        }, stepTime);

        // Flavor Bars
        flavorAcidity.style.width = `${result.flavor.acidity * 10}%`;
        flavorSweetness.style.width = `${result.flavor.sweetness * 10}%`;
        flavorBody.style.width = `${result.flavor.body * 10}%`;
        flavorBitterness.style.width = `${result.flavor.bitterness * 10}%`;

        // Nebula Glow for High Scores
        if (result.total_score >= 95) {
            resultSection.classList.add('nebula-glow');
        } else {
            resultSection.classList.remove('nebula-glow');
        }
    }

    // --- Timer Logic ---
    let timerInterval = null;
    let timerTime = 0; // seconds

    timerToggleBtn.addEventListener('click', () => {
        if (timerInterval) {
            // Stop
            clearInterval(timerInterval);
            timerInterval = null;
            timerToggleBtn.textContent = 'Start';
        } else {
            // Start
            timerInterval = setInterval(() => {
                timerTime++;
                updateTimerDisplay();
            }, 1000);
            timerToggleBtn.textContent = 'Stop';
        }
    });

    timerResetBtn.addEventListener('click', () => {
        clearInterval(timerInterval);
        timerInterval = null;
        timerTime = 0;
        timerToggleBtn.textContent = 'Start';
        updateTimerDisplay();
    });

    function updateTimerDisplay() {
        const m = Math.floor(timerTime / 60);
        const s = timerTime % 60;
        timerMinutes.textContent = m.toString().padStart(2, '0');
        timerSeconds.textContent = s.toString().padStart(2, '0');
    }

    // --- History Logic (Minimal) ---
    function saveToHistory(recipe, result) {
        const history = JSON.parse(localStorage.getItem('brew_history') || '[]');
        const entry = {
            id: Date.now(),
            recipe,
            result,
            date: new Date().toISOString()
        };
        history.unshift(entry);
        localStorage.setItem('brew_history', JSON.stringify(history.slice(0, 10))); // Max 10 items
    }

    openHistoryBtn.addEventListener('click', () => {
        inputSection.classList.add('hidden');
        resultSection.classList.add('hidden');
        historySection.classList.remove('hidden');
        renderHistory();
    });

    historyBackBtn.addEventListener('click', () => {
        historySection.classList.add('hidden');
        inputSection.classList.remove('hidden');
    });

    clearHistoryBtn.addEventListener('click', () => {
        if (confirm('Hapus semua riwayat seduhan?')) {
            localStorage.removeItem('brew_history');
            renderHistory();
        }
    });

    function renderHistory() {
        const history = JSON.parse(localStorage.getItem('brew_history') || '[]');
        historyContainer.innerHTML = '';

        if (history.length === 0) {
            historyContainer.innerHTML = '<p style="text-align:center; color:var(--text-muted); opacity:0.6; margin: 2rem 0;">Belum ada riwayat seduhan.</p>';
            return;
        }

        history.forEach(item => {
            const date = new Date(item.date).toLocaleDateString('id-ID', {
                day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit'
            });
            const div = document.createElement('div');
            div.className = 'history-item';
            div.innerHTML = `
                <div class="history-info">
                    <h4>${item.recipe.method.toUpperCase()}</h4>
                    <p>${date} • ${item.recipe.coffee_dose}g / ${item.recipe.water_yield}g</p>
                </div>
                <div class="history-score">${item.result.total_score}</div>
            `;
            historyContainer.appendChild(div);
        });
    }

    // --- Star Generation ---
    function initStars() {
        const container = document.getElementById('stars-container');
        const count = 50;
        for (let i = 0; i < count; i++) {
            const star = document.createElement('div');
            star.className = 'star';
            const x = Math.random() * 100;
            const y = Math.random() * 100;
            const size = Math.random() * 2 + 1;
            const duration = Math.random() * 3 + 2;
            const delay = Math.random() * 5;

            star.style.left = `${x}%`;
            star.style.top = `${y}%`;
            star.style.width = `${size}px`;
            star.style.height = `${size}px`;
            star.style.setProperty('--duration', `${duration}s`);
            star.style.animationDelay = `${delay}s`;

            container.appendChild(star);
        }
    }

    initStars();

    // --- Phase 3 Logic ---

    // Copy Recipe
    copyRecipeBtn.addEventListener('click', () => {
        const method = document.getElementById('method').value;
        const dose = document.getElementById('coffee_dose').value;
        const yield = document.getElementById('water_yield').value;
        const temp = document.getElementById('temperature').value;
        const time = document.getElementById('brew_time').value;
        const score = totalScoreText.textContent;

        const text = `Celestial Brew Recipe ☕✨\nMethod: ${method.toUpperCase()}\nCoffee: ${dose}g\nWater: ${yield}g\nTemp: ${temp}°C\nTime: ${time}s\nScore: ${score}/100`;

        navigator.clipboard.writeText(text).then(() => {
            const originalText = copyRecipeBtn.innerHTML;
            copyRecipeBtn.textContent = 'Berhasil Disalin!';
            setTimeout(() => {
                copyRecipeBtn.innerHTML = originalText;
            }, 2000);
        });
    });

    // Ratio Calculator
    openCalcBtn.addEventListener('click', () => {
        calcModal.classList.remove('hidden');
        calcCoffeeDose.value = document.getElementById('coffee_dose').value;
    });

    closeCalcBtn.addEventListener('click', () => {
        calcModal.classList.add('hidden');
    });

    const updateCalc = () => {
        const ratio = parseFloat(calcTargetRatio.value) || 0;
        const dose = parseFloat(calcCoffeeDose.value) || 0;
        calcWaterResult.textContent = Math.round(dose * ratio);
    };

    calcTargetRatio.addEventListener('input', updateCalc);
    calcCoffeeDose.addEventListener('input', updateCalc);

    applyCalcBtn.addEventListener('click', () => {
        document.getElementById('coffee_dose').value = calcCoffeeDose.value;
        document.getElementById('water_yield').value = calcWaterResult.textContent;
        calcModal.classList.add('hidden');
    });
});
