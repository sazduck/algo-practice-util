import os
import re
import sys
from pathlib import Path
from string import Template

SCRIPT_DIR = Path(__file__).resolve().parent


def main():
    if len(sys.argv) < 2:
        print("❌ не передано имя модуля", file=sys.stderr)
        return 1
    if len(sys.argv) < 3:
        print("❌ не передано имя папки", file=sys.stderr)
        return 1

    create_next_task(sys.argv[1], sys.argv[2])
    return 0


def create_next_task(platform: str, dirname: str):
    # Регулярка для поиска папок типа task001
    pattern = re.compile(r"^task(\d{3})$")

    # Сканируем текущую директорию
    tasks = [0]  # Начинаем с 0 на случай, если папок еще нет

    mod_path = Path(".", platform).resolve()
    dir_root_path = mod_path / "tasks" / dirname

    os.makedirs(dir_root_path, exist_ok=True)
    for entry in os.listdir(dir_root_path):
        dir_path = dir_root_path / entry
        if os.path.isdir(dir_path):
            match = pattern.match(entry)
            if match:
                tasks.append(int(match.group(1)))

    # Определяем номер следующей задачи
    next_num = max(tasks) + 1
    next_task_name = f"task{next_num:03d}"

    # Создаем папку
    os.makedirs(dir_root_path / next_task_name, exist_ok=True)

    template_path = SCRIPT_DIR / "templates"

    content: str  # Содержимое main.go
    with open(template_path / "main.go.tpl", encoding="utf-8") as f:
        content = f.read()

    # Содержимое main_test.go

    tmpl: Template
    with open(template_path / "main_test.go.tpl", encoding="utf-8") as f:
        tmpl = Template(f.read())

    parts = dir_root_path.parts
    mod_name = parse_go_mod(mod_path) or ""
    try:
        idx = parts.index(platform)
        mod_name = Path(mod_name, *parts[idx:]).as_posix()

    except ValueError:
        print(f"❌ Папка {platform} не найдена в пути {dir_root_path}")

    test_content = tmpl.substitute(
        mod_name=mod_name,
        task_name=next_task_name,
    )

    # Записываем файлы
    files = {
        "solution.go": content,
        "solution_test.go": test_content,
        "test_cases.json": "[]",  # Пустой массив для корректной работы json.Unmarshal
    }

    for filename, content in files.items():
        with open(
            os.path.join(dir_root_path, next_task_name, filename),
            "x",
            encoding="utf-8",
        ) as f:
            f.write(content)

    print(
        f"✅ Создана задача {dir_root_path.relative_to(SCRIPT_DIR, walk_up=True)}:\n"
        + f"   - {next_task_name}/main.go\n"
        + f"   - {next_task_name}/main_test.go\n"
        + f"   - {next_task_name}/test_cases.json\n"
    )


def parse_go_mod(mod_path: Path):
    with open(mod_path / "go.mod", "r", encoding="utf-8") as f:
        content = f.read()
        # Ищем строку, начинающуюся с module, и захватываем следующее слово
        match = re.search(r"^module\s+(.+)$", content, re.MULTILINE)
        if match:
            return match.group(1).strip()
        return ""


if __name__ == "__main__":
    sys.exit(main())
