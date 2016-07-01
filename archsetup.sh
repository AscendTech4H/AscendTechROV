yaourt -S ravi arduino arduino-mk
echo "Now to add you to uucp. . ."
su -c "gpasswd -a "$(whoami)" uucp"
echo "Now to add you to lock. . ."
su -c "gpasswd -a "$(whoami)" lock"
