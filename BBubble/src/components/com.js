document.addEventListener("DOMContentLoaded", function () {
	// 버튼 선택 기능 - 선택된 버튼에 'selected' 클래스를 추가/제거
	document.querySelectorAll('.buttons .button').forEach(button => {
		button.addEventListener('click', () => {
			button.classList.toggle('selected');
		});
	});
});

// 폼 제출 함수
function submitForm() {
	const name = document.getElementById('name').value.trim();
	const age = document.getElementById('age').value.trim();

	// 이름 및 나이 입력 여부 확인
	if (!name) {
		alert("이름을 입력해 주세요.");
		return;
	}
	if (!age) {
		alert("나이를 입력해 주세요.");
		return;
	}

	// 선택된 법률 분야 목록 가져오기
	const selectedLaws = Array.from(document.querySelectorAll('#law-buttons .selected')).map(button => button.innerText);
	if (selectedLaws.length === 0) {
		alert("적어도 하나의 법률 분야를 선택해 주세요.");
		return;
	}

	// 선택된 목소리 가져오기
	const selectedVoice = Array.from(document.querySelectorAll('#voice-buttons .selected')).map(button => button.innerText);
	if (selectedVoice.length === 0) {
		alert("목소리를 선택해 주세요.");
		return;
	}

	// 로딩 화면 표시
	document.getElementById('loading-screen').style.display = 'flex';

	// 로딩 시간을 시뮬레이션하기 위해 setTimeout 사용
	setTimeout(() => {
		// 모달 창에 선택된 정보 표시
		const modalText = `이름: ${name}<br>나이: ${age}<br>선택된 법률 분야: ${selectedLaws.join(", ")}<br>목소리: ${selectedVoice.join(", ")}`;
		document.getElementById('modal-text').innerHTML = modalText;

		// 로딩 화면 숨기기 및 모달 창 열기
		document.getElementById('loading-screen').style.display = 'none';
		document.getElementById('info-modal').style.display = 'flex';
	}, 2000); // 2초 동안 로딩 화면 표시
}


// 모달 닫기 및 페이지 이동 함수
function closeModal() {
	document.getElementById('info-modal').style.display = 'none';
	window.location.href = 'graph.html'; // graph.html로 이동
}



