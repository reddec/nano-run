For Ubuntu/Debian (should be for all LTS)

```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
echo "deb https://dl.bintray.com/reddec/debian all main" | sudo tee /etc/apt/sources.list.d/reddec.list
sudo apt update
sudo apt install nano-run
```

Ansible snippet

```yaml
- name: Add an apt key by id from a keyserver
  become: yes
  apt_key:
    keyserver: keyserver.ubuntu.com
    id: 379CE192D401AB61
- name: Add repository
  become: yes
  apt_repository:
    repo: deb https://dl.bintray.com/reddec/debian all main
    state: present
    filename: reddec
- name: Install nano-run
  become: yes
  apt:
    name: nano-run
    update_cache: yes
    state: latest
```