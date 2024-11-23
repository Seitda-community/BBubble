// 퀴즈 데이터
const quizData = [
    {
        question: "AI 기술을 이용해 얼굴과 목소리를 조작하여 타인의 모습을 흉내내는 기술을 무엇일까요?",
        answer: "딥페이크",
        hint: "ㄷㅍㅇㅋ"
    },
    {
        question: "소비자가 물건을 구입한 후 일정 기간 동안 계약을 철회할 수 있는 권리를 무엇이라 할까요?",
        answer: "청약철회",
        hint: "ㅊㅇㅊㅎ"
    },
    {
        question: "국가가 노사 간의 임금 결정 과정에 개입하여 최저 수준을 정하고 지급을 강제하는 제도를 무엇이라 할까요?",
        answer: "최저임금",
        hint: "ㅊㅈㅇㄱ"
    }
];

let currentQuestionIndex = 0;

// 초기 문제 로드
document.addEventListener("DOMContentLoaded", loadQuestion);

// 문제 로드 함수
function loadQuestion() {
    const questionElement = document.getElementById("question");
    const hintElement = document.getElementById("hint");
    const answerInputContainer = document.getElementById("answerInputContainer");

    // 현재 문제와 초기 힌트 표시
    questionElement.textContent = quizData[currentQuestionIndex].question;
    hintElement.textContent = "힌트를 클릭하세요";

    // 입력 필드 초기화
    answerInputContainer.innerHTML = ""; // 기존 입력 칸을 제거

	// 정답의 글자 수만큼 입력 필드 생성
	const answerLength = quizData[currentQuestionIndex].answer.length;
	for (let i = 0; i < answerLength; i++) {
		const input = document.createElement("input");
		input.type = "text";
		input.classList.add("input-box");
		input.maxLength = 1; // 한 글자만 입력 가능

		// 한 글자가 완성되었을 때 다음 칸으로 포커스 이동
		// input.addEventListener("input", function () {
		// 완성된 한글 또는 한 글자가 입력되었을 때
		// 	if (input.value.length === 1 && i < answerLength - 1) {
		// 		answerInputContainer.children[i + 1].focus();
		// 	}
		// });

		// Backspace 키를 눌렀을 때 이전 칸으로 포커스 이동
		input.addEventListener("keydown", function (event) {
			if (event.key === "Backspace" && input.value === "" && i > 0) {
				answerInputContainer.children[i - 1].focus();
			}
		});

		answerInputContainer.appendChild(input);
	}
}

// 힌트 표시 함수
function showHint() {
    const hintElement = document.getElementById("hint");
    hintElement.innerText = quizData[currentQuestionIndex].hint;
}

// 정답 확인 함수
function checkAnswer() {
    const answerInputContainer = document.getElementById("answerInputContainer");
    let userAnswer = "";

    // 사용자가 입력한 각 문자들을 하나의 문자열로 결합
    for (let i = 0; i < answerInputContainer.children.length; i++) {
        userAnswer += answerInputContainer.children[i].value;
    }

    // 정답과 일치 여부 확인
    if (userAnswer === quizData[currentQuestionIndex].answer) {
        alert("정답입니다!");
    } else {
        alert("오답입니다. 다시 시도해 보세요!");
    }
}

// 정답 모달 표시 함수
function showAnswer() {
    const answerText = document.getElementById("answerText");
    answerText.textContent = quizData[currentQuestionIndex].answer;
    document.getElementById("answerModal").style.display = "flex";
}

// 모달 닫기 함수
function closeModal() {
    document.getElementById("answerModal").style.display = "none";
}

// 다음 문제로 이동
function nextQuestion() {
    currentQuestionIndex = (currentQuestionIndex + 1) % quizData.length;
    loadQuestion();
}

// 모달 바깥을 클릭하면 모달을 닫기
window.onclick = function(event) {
    const answerModal = document.getElementById("answerModal");
    if (event.target === answerModal) {
        closeModal();
    }
}
