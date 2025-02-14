# gignr

`gignr` is a powerful **CLI tool** designed to help developers fetch, manage, and create `.gitignore` templates with ease. It integrates templates from **GitHub**, **TopTal**, and user-defined repositories.

## ✨ Features

- 📦 **Fetch `.gitignore` templates** from:
  - `GitHub`  
  - `GitHub Global`  
  - `GitHub Community`  
  - `TopTal`  
  - Custom **user-defined repositories**
- 🔍 **TUI-powered template search** with filtering and selection.
- ⚡ **Merge multiple templates** into a single `.gitignore` file.
- 💾 **Save and manage `.gitignore` templates locally**.
- 🛠️ **Highly configurable** via `config.yaml`.

## 📥 Installation

- Using Go 🐹

    ```sh
    go install github.com/jasonuc/gignr@latest
    ```

- Using Homebrew 🍺

    ```sh
    brew tap jasonuc/tap && brew install gignr
    ```

## 📌 Usage

### 🛠️ **Creating a `.gitignore` File**

```sh
gignr create gh:Go tt:clion best_go_gitignore
```

- `gh:` → Fetch from **GitHub**
- `ghg:` → Fetch from **GitHub Global**
- `ghc:` → Fetch from **GitHub Community**
- `tt:` → Fetch from **TopTal**
- *(No prefix)* → Fetch from **locally saved templates**

### 🎯 **Adding a Custom Repository**

```sh
gignr add https://github.com/jasonuc/gitignore-templates -n jc
```

- `-n myrepo` sets a **nickname** for the repository.

### 🔍 **Searching for Templates (TUI)**

```sh
gignr search
```

- **Navigate sources**: `←/→`, `tab`
- **Select template**: `Enter`
- **Filter templates**: Start typing
- **Get command to generate**: `Shift + C`
- **Generate from selection**: `Shift + S`
- **Exit**: `Ctrl + C`

### 💾 **Saving a Custom `.gitignore`**

```sh
gignr save best_go_gitignore
```

- Saves `.gitignore` from the **current directory** to **local storage**.
- Storage path is configurable in `config.yaml`.

## ⚙️ Configuration (`config.yaml`)

Located at: `~/.config/gignr/config.yaml`

```yaml
templates:
  storage_path: "~/.config/gignr/templates"
repositories:
  jc: "https://github.com/jasonuc/gitignore-templates"
```

## 🤝 Contributing

Contributions are welcome!  
Fork the repo, make your changes, and open a **Pull Request** 🚀

## 📜 License

This project is licensed under the **MIT License**.
