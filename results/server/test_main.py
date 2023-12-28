from server.tag import Tag
import datetime
import pytest
import sys

t = Tag("PV001", 10)
print(t)


def test_less_than():
    print(f"Tag Value: {t.value}")
    print(f"Tag Less Than: {t.less_than}")
    assert t.test_less_than() == True


def test_greater_than():
    print(f"Tag Value: {t.value}")
    print(f"Tag Greater Than: {t.greater_than}")
    assert t.test_greater_than() == True


if __name__ == "__main__":
    dt = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
    html = f"--html=logs/pytest_report_{dt}.html"
    pytest.main(
        [
            "-v",
            "--showlocals",
            "--durations=10",
            "server/test_main.py",
            html,
            "--self-contained-html",
        ]
    )
