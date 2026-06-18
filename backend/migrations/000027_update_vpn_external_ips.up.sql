-- Update VPN URLs to use external IPs (accessible from production server)
UPDATE external_services SET name='VPN DevOps', url='https://35.219.27.243' WHERE url='https://10.184.20.11';
UPDATE external_services SET name='VPN IT', url='https://35.219.125.94' WHERE url='https://10.184.20.13';
UPDATE external_services SET name='VPN IT DevOps', url='https://34.101.125.72' WHERE url='https://10.88.12.10';
UPDATE external_services SET name='VPN IT RnD', url='https://34.128.98.215' WHERE url='https://10.184.20.14';
