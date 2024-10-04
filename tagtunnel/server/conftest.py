""" conf for reports """

import pytest
from datetime import datetime

IMG = "assets/fe_logo.png"


def pytest_html_report_title(report):
    report.title = "MACFE HIL Report"


pytest.hookimpl(optionalhook=True)


def pytest_html_results_summary(prefix, summary, postfix):
    img_html = f'<img src="{IMG}" alt="FE Logo" width="100" height="100">'
    postfix.extend([f"<h1>{img_html}</h1>"])
    summary.extend(
        [
            f"<h2>Hardware-in-the-loop test results of the electrical-software test-bench</h2>"
        ]
    )


def pytest_html_results_table_header(cells):
    cells.pop()
    cells.insert(2, "<th>Condition</th>")
    cells.insert(1, '<th class="sortable time" data-column-type="time">Time</th>')


def pytest_html_results_table_row(report, cells):
    cells.pop()
    cells.insert(2, f"<td>{report.filename}</td>")
    cells.insert(1, f'<td class="col-time">{datetime.utcnow()}</td>')


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item, call):
    """add description to report from docstring"""
    outcome = yield
    report = outcome.get_result()
    report.filename = str(item.function.__doc__)
