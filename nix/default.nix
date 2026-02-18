# vim: ts=2 sw=2
{ fetchFromGitHub
, buildGoModule
, lib
, ...
}: buildGoModule rec {
  pname = "program-browser";
  version = "0.0.5";

  src = fetchFromGitHub {
    owner = "m1dugh";
    repo = "program-browser";
    rev = "v${version}";
    hash = "sha256-UMy+m7V1nGPY8H+mus3ZkxkrWWg/jIpYZ5dF0iB0SZM=";
  };
  vendorHash = "sha256-Afw2gCq4hLv8FG6sBTF0QtymxMOYXPVrawF06UJLUBs=";

  meta = with lib; {
    description = "A tool to list bug bounty programs on plateforms";
    homepage = "https://github.com/m1dugh/program-browser";
    license = licenses.mit;
    maintainers = with maintainers; [ m1dugh ];
  };
}
