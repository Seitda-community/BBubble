# BBubble


1. 프로젝트 개요<br>
버블버블은 네이버 클라우드 챗봇과 AI를 활용하여 신종 범죄 동향과 사용자 기반 맞춤형 정보를 제공하는 서비스입니다. <br>
특히 실시간 키워드 분석을 통해 관련된 주요 키워드를 추천하고, 24시간 자동화된 복지 서비스, 맞춤형 학습 기능, <br>
통합 생성형 지식 서비스(GPT)를 제공합니다.<br><br>


2. 주요 기능<br>
신종 범죄 사용자 정보 제공: 사용자 정보를 바탕으로 최신 범죄 동향과 관련 정보를 제공합니다.<br>
실시간 키워드 분석 및 추천: 실시간 키워드 분석을 통해 사용자에게 맞춤형 키워드를 추천하여 최신 정보를 빠르게 전달합니다.<br>
복지향상 챗봇 (버블버블 챗봇): 사용자 정보 기반으로 24시간 서비스 자동화를 지원하여 사용자의 편의성을 높입니다.<br><br>
개인 맞춤형 학습 서비스 (버블버블 퀴즈): 사용자 정보에 맞춘 개인화된 학습 퀴즈 기능을 제공합니다.<br>
통합 생성형 지식 서비스 (버블버블 GPT): 사용자 정보 기반으로 지식과 정보를 통합하여 제공하는 AI 지식 서비스입니다.<br>


3. 한계점<br>
사용자 정보 활용에 대한 제한: 정보 제공과 맞춤형 서비스는 사용자 정보의 정확성과 완전성에 따라 제한될 수 있습니다.<br>
실시간 서비스 제공 한계: 24시간 자동화된 서비스를 위해 안정적인 시스템 환경과 지속적인 모니터링이 필요합니다.<br><br>


4.시작하기<br>
이 프로젝트를 시작하려면 먼저 이 저장소를 클론하고 아래 절차에 따라 설정합니다.<br><br>

1) 저장소 클론<br>
git clone https://github.com/your-username/your-repo.git<br>
cd your-repo<br><br>

2) 의존성 설치: 프로젝트 디렉터리에서 의존성을 설치합니다.<br>
flutter pub get<br><br>

3) 네이버 클라우드 API 키 설정:<br>
.env 파일에 네이버 클라우드 API 인증 정보를 추가하여 환경 변수를 설정합니다.<br><br>

5. 폴더구조<br> 
project_root/<br>
├── lib/<br>
│   ├── main.dart                 # Flutter 앱의 시작 파일<br>
│   ├── screens/<br>
│   │   └── chat_screen.dart       # 사용자-챗봇 대화 화면<br>
│   ├── services/<br>
│   │   └── chatbot_service.dart   # 네이버 클라우드 챗봇 API 연동 로직<br>
│   └── widgets/<br>
│       └── message_bubble.dart    # 대화 메시지를 표시하는 위젯<br>
├── assets/                        # 이미지, 폰트 등의 자산 파일<br>
├── .env                           # 환경 변수 파일 (민감 정보는 여기에 저장)<br>
├── pubspec.yaml                   # 프로젝트 의존성 파일<br>
└── README.md                      # 프로젝트 설명 파일<br><br>


6. 환경설정<br>
1) API 키 설정:<br>
.env 파일을 생성하고 다음과 같이 API 키와 시크릿 키를 입력합니다.<br>
NAVER_CLOUD_CLIENT_ID=your_client_id<br>
NAVER_CLOUD_CLIENT_SECRET=your_client_secret<br><br>

2) 네이버 클라우드 API 연동:<br>
chatbot_service.dart 파일에 네이버 클라우드 챗봇 API를 호출하는 코드를 설정합니다.<br><br>

7. 사용법<br>
1) 앱 실행: 에뮬레이터 또는 실제 기기에서 앱을 실행합니다.<br><br>
flutter run<br><br>

2) 챗봇과 대화하기: 대화 화면에서 실시간 키워드 추천, 복지 서비스, 개인 맞춤형 퀴즈, 생성형 지식 서비스를 경험할 수 있습니다.<br><br>

8. 라이선스<br>
이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 LICENSE 파일을 참조하세요.<br><br>