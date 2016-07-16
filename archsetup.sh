su -s "pacman -S binutils make gcc fakeroot"
yaourt -S ravi arduino
echo "Now to add you to uucp. . ."
su -c "gpasswd -a "$(whoami)" uucp"
echo "Now to add you to lock. . ."
su -c "gpasswd -a "$(whoami)" lock"
