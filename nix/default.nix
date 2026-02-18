# vim: ts=2 sw=2
{ fetchFromGitHub
, buildGoModule
, lib
, ...
}: buildGoModule rec {
  pname = "program-browser";
  version = "0.0.4";

  src = fetchFromGitHub {
    owner = "m1dugh";
    repo = "program-browser";
    rev = "v${version}";
    hash = "sha256-eSfXgPHyuQ48gu+RmAcNtBiqvVw3wX3s7sGMxEa5Q1w=";
  };
  vendorHash = "sha256-Afw2gCq4hLv8FG6sBTF0QtymxMOYXPVrawF06UJLUBs=";

  meta = with lib; {
    description = "A tool to list bug bounty programs on plateforms";
    homepage = "https://github.com/m1dugh/program-browser";
    license = licenses.mit;
    maintainers = with maintainers; [ m1dugh ];
  };
}
